package nethandler

import (
    "../localip"
    "fmt"
    "os"
    "sync"
    "../peers"
    ."../../types"
    "strings"    
    "../../misc"
    "time"
)

var _mtx sync.Mutex


/********************
* New local elevator
*********************/

func Create_new_local_elevator  (   network_name string,
                                    list_of_elevators * []Elevator,
){
    e := Elevator{}
    e.Id = Set_eID(network_name)
    e.Status = "idle"
    e.Direction = "stopped"

    // add the elevator to the elevator list
    _mtx.Lock()
    (*list_of_elevators) = append((*list_of_elevators), e)
    _mtx.Unlock()
}


/********************
* Create elevatorID
*********************/

func Set_eID(name_of_network string) string{
    
    localIP, err := localip.LocalIP()

    if err != nil {
        fmt.Println(err)
        localIP = "DISCONNECTED"
    }

    return fmt.Sprintf("%s::%d #%s", localIP, os.Getpid(), name_of_network)
}


/******************************************
* The network peer handler will keep track 
* of network connections and then update 
* the elevator connection list
*******************************************/

func Peer_update_handler    (   peer_ch <-chan peers.PeerUpdate, 
                                failed_order_ch chan<- Order,
                                network_name string,
                                status_msg chan<- string,
                                list_of_elevators * []Elevator,
){
    for {       
        select {

        // if any new or lost connections
        case p := <-peer_ch:

            _mtx.Lock()
            // copy an instace of the elevators list
            elevators := (*list_of_elevators)
            _mtx.Unlock()


            // add new connected elevator to elevator list
            if p.New != "" && p.New != elevators[0].Id && strings.Contains(p.New, network_name) {
                e := Elevator{}
                e.Id = p.New
                e.Status = "unknown"
                e.Direction = "unknown"

                _mtx.Lock()
                (*list_of_elevators) = append((*list_of_elevators), e) 
                _mtx.Unlock()

                status_msg <- "Got new elevator, say hello to " + e.Id


            // delete lost elevators if disconnected, but also redistribute any unfinished orders
            } else if p.Lost != nil {
                
                num_of_elevators := len(elevators)
                e_id := elevators[0].Id

                for i:=0; i<num_of_elevators; i++ {
                    for j:=0; j<len(p.Lost); j++ {

                        if elevators[i].Id == p.Lost[j] && p.Lost[j] != e_id && strings.Contains(p.Lost[j], network_name){ 

                            status_msg <- "Peer disconnected, say goodbye to " + elevators[i].Id

                            // if elevator has unfinished orders
                            if elevators[i].Order_list != nil {

                                for o:=0; o<len(elevators[i].Order_list); o++ {     

                                    // return unfinished hall orders to network
                                    if elevators[i].Order_list[o].Button != 2 {
                                        failed_order_ch <- elevators[i].Order_list[o]
                                    }
                                }
                            }
                            // then delete the elevator from the list 
                            _mtx.Lock()
                            (*list_of_elevators) = misc.Remove_elevator(i, (*list_of_elevators))
                            _mtx.Unlock()
                        }
                    }
                }
                
            }
            // update number of connections
            _mtx.Lock()            
            (*list_of_elevators)[0].Connections = len(p.Peers)
            _mtx.Unlock()
        }
    }
}


/***************************************
* The network order receiver will check
* and add incoming network orders   
****************************************/

func Network_order_rx   (   order_buffer_ch chan<- Order, 
                            order_rx_ch chan Order,
                            status_msg chan<- string,
                            list_of_elevators * []Elevator,
){
    for {
        select {

        // if new incoming order
        case incoming_order := <- order_rx_ch:  

             // add to order buffer list if order is assigned to this elevator
            _mtx.Lock()
            if incoming_order.Assigned_to == (*list_of_elevators)[0].Id {        

                // open door if on same floor as order
                if incoming_order.Target_floor == (*list_of_elevators)[0].Current_floor {
                    misc.Open_door(&(*list_of_elevators))
                }

                // if not duplicate order, add to order buffer
                if misc.Is_duplicate_order((*list_of_elevators), incoming_order) == false && incoming_order.Target_floor != (*list_of_elevators)[0].Current_floor{
                    order_buffer_ch <- incoming_order
                    //fmt.Printf("Received new network order!\n")       
                    status_msg <- "Received new network order!"           
                } else {
                    //fmt.Printf("Received duplicate network order...\n")     
                    status_msg <- "Received duplicate network order..."  
                }
            }
            _mtx.Unlock()
        }
    }
}


/*******************************
* This will transmit the local 
* elevator state to the network   
********************************/

func Elevator_state_tx  (   repeat_timer_ms int, 
                            elevator_tx_ch chan Elevator,
                            list_of_elevators * []Elevator,
){ 
    for {
        _mtx.Lock()
        elevator_tx_ch <- (*list_of_elevators)[0]
        _mtx.Unlock()

        time.Sleep(time.Duration(repeat_timer_ms) * time.Millisecond)
    }
}


/*******************************
* This will receive elevator 
* states from the network   
********************************/

func Elevator_state_rx  (   elevator_rx_ch chan Elevator,
                            list_of_elevators * []Elevator,
){ 
    for {
        select {

        // if new elevator state msg
        case e := <- elevator_rx_ch:    

            // make local copy of elevators list
            _mtx.Lock()
            elevators := (*list_of_elevators)
            _mtx.Unlock()

            num_of_elevators := len(elevators)

            // check if this is a known elevator
            for i:=0; i<num_of_elevators; i++ {   
                if elevators[i].Id == e.Id {
                    // update known elevator state
                    _mtx.Lock()
                    (*list_of_elevators)[i] = e
                    _mtx.Unlock()
                }
            }
        }
    }
}