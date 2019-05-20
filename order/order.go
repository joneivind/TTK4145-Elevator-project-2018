package order

import  (
    ."../types"
    "../misc"    
    "../elevator_io"
    "fmt"
    "math/rand"
    "encoding/json"
    "sync"
    "strings"
    "io/ioutil"
    "strconv"
    "time"
)

var _mtx sync.Mutex


/*********************************
* Write cab orders to file
**********************************/

func Write_cab_orders_to_file   (   list_of_elevators * []Elevator,
){

    var stop_on_floors []int

    // search for cab orders on local elevator
    _mtx.Lock()
    for c:=0; c<len((*list_of_elevators)[0].Order_list); c++{
        if (*list_of_elevators)[0].Order_list[c].Button == 2 {            
            stop_on_floors = append(stop_on_floors, (*list_of_elevators)[0].Order_list[c].Target_floor)   
        }
    }
    _mtx.Unlock()

    // write cab orders to file
    rankingsJson, _ := json.Marshal(stop_on_floors)
    err := ioutil.WriteFile("cab_orders.json", rankingsJson, 0644)
    if err != nil {
        fmt.Println("File error: %v\n\n", err)
    }
}


/*********************************
* Read cab orders from file
**********************************/

func Add_cab_orders_from_file   (   order_buffer_ch chan Order, 
                                    num_of_floors int, 
                                    status_msg chan<- string,
                                    list_of_elevators * []Elevator,
){
    orders_found := false

    // read cab orders from file
    file, err := ioutil.ReadFile("cab_orders.json")
    if err != nil {
        fmt.Printf("File error: %v\n", err)
    }

    // send all backup cab orders to order buffer
    for f:=0; f<num_of_floors; f++ {
        if strings.Contains(string(file), strconv.Itoa(f)) { 

            orders_found = true

            status_msg <- "Found backup cab orders"

            cab_order := Order{}
            cab_order.Id = rand.Intn(10000)

            _mtx.Lock()
            cab_order.Origin_e = (*list_of_elevators)[0].Id
            cab_order.Assigned_to = (*list_of_elevators)[0].Id
            _mtx.Unlock()

            cab_order.Button = 2
            cab_order.Target_floor = f

            order_buffer_ch <- cab_order
        }
    }
    if orders_found != true { status_msg <- "No previous cab orders found..." }
}



/*********************************************
* The Order assigner will receive input from
* the button panel or any local failed order
* and assign a host based on cost of the order
**********************************************/

func Order_assigner (   drv_buttons <-chan elevator_io.ButtonEvent, 
                        order_buffer_ch chan<- Order, 
                        order_tx_ch chan<- Order, 
                        failed_order_ch chan Order,
                        status_msg chan<- string,
                        list_of_elevators * []Elevator,
){
    
    var button_pressed elevator_io.ButtonEvent
    rand.Seed(42)


    for {
        select {

        // if elevator button pressed
        case button_pressed = <- drv_buttons:


            // create new temp order
            incoming_order := Order{}
            incoming_order.Id = rand.Intn(10000)
            _mtx.Lock()
            incoming_order.Origin_e = (*list_of_elevators)[0].Id
            _mtx.Unlock()
            incoming_order.Button = button_pressed.Button
            incoming_order.Target_floor = button_pressed.Floor
            


            // if cab order
            if incoming_order.Button == 2 {

                // if not duplicate order, add to order buffer  
                if misc.Is_duplicate_order((*list_of_elevators), incoming_order) == false {          
                    incoming_order.Assigned_to = (*list_of_elevators)[0].Id
                    order_buffer_ch <- incoming_order   // send cab order to this elevator
                    status_msg <- "Received new local cab order!"
                
                } else {
                    status_msg <- "Received duplicate local cab order..."
                }    



            // if hall order
            } else {    

                var cost_list []Cost
                var lowest_cost Cost

                _mtx.Lock()
                num_of_elevators := len((*list_of_elevators))
                _mtx.Unlock()
                
                lowest_cost.Sum = 100000


                // warning msg when only one elevator on network
                _mtx.Lock()
                if (*list_of_elevators)[0].Connections <= 1 { status_msg <- "Cannot currently accept hall orders..." }
                _mtx.Unlock()

                // check order cost for every elevator
                _mtx.Lock()
                for i:=0; i<num_of_elevators; i++ {     
                   cost_list = append(cost_list, misc.Order_cost(incoming_order, (*list_of_elevators)[i]))
                }
                _mtx.Unlock()

                // find the elevator with the lowest cost
                for i:=0; i<len(cost_list); i++ {     
                    if cost_list[i].Sum < lowest_cost.Sum {     
                        lowest_cost = cost_list[i]                            
                    }
                }

                // set assigned order id to the elevator with lowest cost
                incoming_order.Assigned_to = lowest_cost.Elevator_id    


                // send order to network if not assigned to self or if error
                if incoming_order.Assigned_to != (*list_of_elevators)[0].Id || (*list_of_elevators)[0].Status == "error" {    

                    // check for duplicate order
                    if misc.Is_duplicate_order((*list_of_elevators), incoming_order) == false {    

                        // if not duplicate order, send hall order to network
                        order_tx_ch <- incoming_order    
                        status_msg <- "Hall order assigned to " + incoming_order.Assigned_to

                    } else { status_msg <- "Received duplicate local hall order..." }
                    


                // only accept local hall orders if more than one elevator on the network
                } else if incoming_order.Assigned_to == (*list_of_elevators)[0].Id && (*list_of_elevators)[0].Connections > 1 {
                    
                    // if not duplicate order, add to internal order buffer
                    if misc.Is_duplicate_order((*list_of_elevators), incoming_order) == false {
                        order_buffer_ch <- incoming_order
                        status_msg <- "Received new local hall order!"
                    } else { status_msg <- "Received duplicate local hall order..." }
                }
            }
                


        // failed order redistribution     
        case failed := <- failed_order_ch:     

            var cost_list []Cost
            var lowest_cost Cost


            // reassign failed order
            incoming_order := failed

            // tag failed order with error code
            incoming_order.Error = 1    

            lowest_cost.Sum = 100000
            


            // checking order cost for every elevator
            for i:=0; i<len((*list_of_elevators)); i++ {     
               cost_list = append(cost_list, misc.Order_cost(incoming_order, (*list_of_elevators)[i]))
            }
            // find the elevator with the lowest cost
            for i:=0; i<len(cost_list); i++ {     
                if cost_list[i].Sum < lowest_cost.Sum {     
                    lowest_cost = cost_list[i]                            
                }
            }
            // assign new elevator to order
            incoming_order.Assigned_to = lowest_cost.Elevator_id    



            // send failed hall order to network
            if incoming_order.Assigned_to != (*list_of_elevators)[0].Id || (*list_of_elevators)[0].Status == "error" {
                status_msg <- "Failed hall order is reassigned to: " + incoming_order.Assigned_to
                order_tx_ch <- incoming_order    
            
            // or to self if it has lowest cost
            } else if incoming_order.Assigned_to == (*list_of_elevators)[0].Id && (*list_of_elevators)[0].Connections > 0 && incoming_order.Target_floor != (*list_of_elevators)[0].Current_floor {
        
                // if not duplicate order, add to local order buffer
                if misc.Is_duplicate_order((*list_of_elevators), incoming_order) == false {     
                    order_buffer_ch <- incoming_order
                    status_msg <- "Received new local hall order!" 

                } else { status_msg <- "Received duplicate local hall order..." }
            }
        }
    }
}


/******************************************
* The order handler will read incoming
* orders from buffer and send them to the 
* elevator. Then starts a timer for the order 
* list. Deletes order when ACK msg is 
* received from elevator, or sending failed
* order back to order assigner
*******************************************/

func Order_handler  (   order_buffer_ch <-chan Order, 
                        orders_ch chan<- Order, 
                        order_ack_ch <-chan int, 
                        order_error_ch <-chan int, 
                        order_tx_ch chan<- Order, 
                        failed_order_ch chan Order,
                        status_msg chan<- string,
                        list_of_elevators * []Elevator,
){


    var order_list []Order // store all incoming orders
    var failed_o Order // store all failed orders

    // create new watchdog channel
    watchdog_ack_timer_ch := make(<-chan time.Time)    
    watchdog_timeout_sec := 30 * time.Second


    for {
        select {


        // new incoming order from order assigner
        case incoming_order := <- order_buffer_ch: 

                // create new order
                new_order := Order{}
                new_order.Id = incoming_order.Id
                new_order.Origin_e = incoming_order.Origin_e
                new_order.Assigned_to = incoming_order.Assigned_to
                new_order.Button = incoming_order.Button
                new_order.Target_floor = incoming_order.Target_floor
                
                // add new order to order list
                order_list = misc.Push_order(new_order, order_list)

                // start watchdog timer
                watchdog_ack_timer_ch = time.After(watchdog_timeout_sec)



        // ACK msg from elevator
        case ack_id := <- order_ack_ch:   

            elevators := (*list_of_elevators)
            
            // restart new watchdog timer
            watchdog_ack_timer_ch = time.After(watchdog_timeout_sec)

            // check for matching ACK id in order list
            for i:=0; i<len(order_list); i++ {

                // if a matching order id is found
                if order_list[i].Id == ack_id {

                    // change ACK status on order
                    //status_msg <- "ACK received for order " + string(order_list[i].Id)
                    order_list[i].Ack_received = true
                    
                    // remove finished order from order list
                    order_list = misc.Remove_order(i, order_list)

                    // update elevator order list
                    if elevators[0].Order_list != nil {
                        (*list_of_elevators)[0].Order_list = misc.Remove_order(i, elevators[0].Order_list)
                    }

                    // if order list is empty, change status to idle
                    if order_list == nil {  
                        (*list_of_elevators)[0].Status = "idle"
                        status_msg <- "No orders..."
                    }
                }
            }
            // update cab order file backup
            go Write_cab_orders_to_file(&elevators)



        // if order has failed to complete
        case failed_order_id := <- order_error_ch:

            // check for matching id in order list
            for i:=0; i<len(order_list); i++ {

                // if a matching order id is found
                if order_list[i].Id == failed_order_id {

                    order_list[i].Assigned_to = ""
                    order_list[i].Error = 1
                    
                    status_msg <- "Got order timeout error"

                    failed_o = order_list[i]
                    
                    // update elevator order list
                    (*list_of_elevators)[0].Order_list = misc.Remove_order(i, (*list_of_elevators)[0].Order_list)

                    // remove order from order list
                    order_list = misc.Remove_order(i, order_list)

                    // send failed hall order back to order assigner
                    if failed_o.Button != 2 {
                        failed_order_ch <- failed_o
                    }
                }
            }
            // update cab order file backup
            go Write_cab_orders_to_file(&(*list_of_elevators))


        // if elevator has not received ACK for all orders within timeout
        case <- watchdog_ack_timer_ch:
        
            if order_list != nil {

                (*list_of_elevators)[0].Status = "error"

                // check for matching id
                for i:=0; i<len(order_list); i++ {  

                    order_list[i].Assigned_to = ""
                    order_list[i].Error = 1
                    
                    status_msg <- "ACK timed out, returning unfinished orders to order handler"

                    // store failed order
                    failed_o = order_list[i]    
                    
                    // update elevator order list
                    (*list_of_elevators)[0].Order_list = misc.Remove_order(i, (*list_of_elevators)[0].Order_list)   

                    // remove order from order list
                    order_list = misc.Remove_order(i, order_list)    

                    // send failed hall order back to order assigner
                    if failed_o.Button != 2 {      
                        failed_order_ch <- failed_o
                    }
                }
            }
            // update cab order file backup
            go Write_cab_orders_to_file(&(*list_of_elevators))


        default:
            
            // run through order list
            if order_list != nil {
                
                // send all new orders to elevator
                for i:=0; i<len(order_list); i++ {
                    
                    // check if order is already sent
                    if order_list[i].Is_sent == false {

                        // send order
                        orders_ch <- order_list[i]
                        order_list[i].Is_sent = true

                        // update elevator order list
                        (*list_of_elevators)[0].Order_list = misc.Push_order(order_list[i], (*list_of_elevators)[0].Order_list)
                        (*list_of_elevators)[0].Status = "busy"

                        // update cab order file backup
                        if order_list[i].Button == 2 { go Write_cab_orders_to_file(&(*list_of_elevators)) }                        
                    }
                }
            }
            time.Sleep(10 * time.Millisecond) 
        }
    }
}


/******************************************    
* This function will receive orders on the 
* orders channel and store them in a buffer.
* The orders will be executed depending on
* current direction and an ACK is then
* published on the order_complete channel
* when an order is complete. Has timeout if
* floor is not reached in time.
*******************************************/

func Order_executer     (   orders_ch <-chan Order, 
                            order_ack_ch chan<- int, 
                            drv_floors <-chan int, 
                            drv_stop <-chan bool, 
                            order_error_ch chan<- int,
                            status_msg chan<- string,
                            list_of_elevators * []Elevator,
){

    var order_buffer []Order    // buffer incoming orders
    current_floor := 0          // store current floor
    stop_flag := false
    
    // create new watchdog channel
    watchdog_timer := make(<-chan time.Time)    
    watchdog_timeout_sec := 4 * time.Second 

    for{
        select{

        // when new orders are received
        case new_order := <- orders_ch: 

            // open door if on same floor
            if new_order.Target_floor == (*list_of_elevators)[0].Current_floor {
                misc.Open_door(&(*list_of_elevators))
                order_ack_ch <- new_order.Id
            
            } else {                    
                order_buffer = misc.Push_order(new_order, order_buffer)   
                // start timer      
                watchdog_timer = time.After(watchdog_timeout_sec)      
            }


        // new floor
        case current_floor = <- drv_floors:     

            // zero out timer
            watchdog_timer = nil    

            elevator_io.SetFloorIndicator(current_floor)
            _mtx.Lock()
            // update floor status
            (*list_of_elevators)[0].Current_floor = current_floor     
            _mtx.Unlock()

            // stop elevator when limit of floors is reached  
            if current_floor >= 3 || current_floor <= 0 || stop_flag == true {                 
                elevator_io.SetMotorDirection(elevator_io.MD_Stop)                    
                (*list_of_elevators)[0].Direction = "stopped"
            }

            // stop elevator when stop button is pressed   
            if stop_flag == true {                
                elevator_io.SetMotorDirection(elevator_io.MD_Stop)                    
                (*list_of_elevators)[0].Direction = "stopped"
                stop_flag = true

                misc.Open_door(&(*list_of_elevators))

                stop_flag = false
                (*list_of_elevators)[0].Status = "busy"   
                elevator_io.SetStopLamp(false)    
            }

            // search through order buffer for matching orders when on new floor
            for i := 0; i < len(order_buffer); i++ {    

                // if an order is found on this floor
                if order_buffer[i].Target_floor == current_floor {   

                    elevator_io.SetMotorDirection(elevator_io.MD_Stop)
                    (*list_of_elevators)[0].Direction = "stopped"

                    order_ack_ch <- order_buffer[i].Id 

                    _mtx.Lock()
                    order_buffer = misc.Remove_order(i, order_buffer)
                    _mtx.Unlock()

                    misc.Open_door(&(*list_of_elevators))
                }
            }


        // if stop button is pressed, raise flag
        case stop_button_pressed := <-drv_stop:

            if stop_button_pressed {
                stop_flag = true
                elevator_io.SetStopLamp(true)              
                (*list_of_elevators)[0].Status = "Stop"
                status_msg <- "Stop button is pressed, will stop at next floor"
            }


        // if elevator cannot fulfill order in time
        case <- watchdog_timer:    

            if order_buffer != nil {

                elevator_io.SetMotorDirection(elevator_io.MD_Stop)
                (*list_of_elevators)[0].Status = "error"

                status_msg <- "An order has timed out, returning unfinished orders to order handler"
                
                for i := 0; i < len(order_buffer); i++ {
                    order_error_ch <- order_buffer[i].Id

                    _mtx.Lock()
                    order_buffer = misc.Remove_order(i, order_buffer)
                    _mtx.Unlock()
                }
            }


        default:
            
            // check if any orders in order buffer
            if order_buffer != nil {    

                // if order is on same floor
                if current_floor == order_buffer[0].Target_floor {
                    
                    elevator_io.SetMotorDirection(elevator_io.MD_Stop)

                    order_ack_ch <- order_buffer[0].Id 
                    order_buffer = misc.Remove_order(0, order_buffer)
                
                // if order floor is above
                } else if current_floor < order_buffer[0].Target_floor {      

                    elevator_io.SetMotorDirection(elevator_io.MD_Up) 
                    (*list_of_elevators)[0].Direction = "up"     

                    if watchdog_timer == nil {
                        watchdog_timer = time.After(watchdog_timeout_sec) 
                    }

                // if order floor is below
                } else if current_floor > order_buffer[0].Target_floor { 

                    elevator_io.SetMotorDirection(elevator_io.MD_Down)    
                    (*list_of_elevators)[0].Direction = "down"
                    
                    if watchdog_timer == nil {
                        watchdog_timer = time.After(watchdog_timeout_sec) 
                    }
                }
            }
            time.Sleep(20 * time.Millisecond) 
        }
    }
}
