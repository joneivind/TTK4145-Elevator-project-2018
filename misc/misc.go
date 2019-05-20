package misc

import (
    ."../types"
    "math"
    "sync"
    "fmt"
    "time"
    "os"
    "os/exec"
    "../elevator_io"
    "encoding/json"
)

var _mtx sync.Mutex



/*********************************
* This is the terminal interface
**********************************/

func Print_state    (   status_msg <-chan string,
                        drv_obstr chan bool,
                        list_of_elevators * []Elevator,
){  
    var new_status_msg string
    var animate bool
    swap_counter := 0

    for {
        select{
        case new_status_msg = <- status_msg:
        case animate = <- drv_obstr:
        default:

            c := exec.Command("clear")
            c.Stdout = os.Stdout
            c.Run()

            fmt.Println(" ")
            fmt.Println(" ")
            fmt.Println(" ")
           
            if swap_counter < 2 || (swap_counter >= 4 && swap_counter < 6)  {                
                fmt.Println("      ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ ")
                fmt.Println("     ░░░▀░░░░░▄▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▄░░░░░▀░░")
                fmt.Println("     ░░░░░░░░█░░▄▀▀▀▀▀▀▀▀▀▀▀▀▀▄░░█░░░░░░░")
                fmt.Println("     ░░░░░░░░█░█░▀░░░░░▀░░▀░░░░█░█░░░░░░░")
                fmt.Println("     ░░░░░░░░█░█░░░░░░░░▄▀▀▄░▀░█░█▄▀▀▄░░░")
                fmt.Println("     ░░▄▀▀█▄░█░█░░▀░░░░░█░░░▀▄▄█▄▀░░░█░░░")
                fmt.Println("     ░░▀▄▄░▀██░█▄░▀░░░▄▄▀░░░░░░░░░░░░▀▄░░")
                fmt.Println("     ░░░░▀█▄▄█░█░░░░▄░░█░░░▄█░░░▄░▄█░░█░░")
                fmt.Println("     ░░░░░░░▀█░▀▄▀░░░░░█░██░▄░░▄░░▄░███░░")
                fmt.Println("     ░░░░░░░░█▄░░▀▀▀▀▀▀▀▀▄░░▀▀▀▀▀▀▀░▄▀░░░")
                fmt.Printf("     ░░░░░░█░░▄█▀█▀▀█▀▀▀▀▀▀▀█▀▀█▀█▀▀█░░░░   - You are elevator #1 (%v)\n", (*list_of_elevators)[0].Id)
                fmt.Println("     ░░▀░░▀▀▀▀░░▀▀▀░░░░░░░░░▀▀▀░░▀▀░░░▀░░░")

            } else if swap_counter >= 6 && swap_counter <= 8 {
                fmt.Println("      ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ ")
                fmt.Println("     ░░░▀░░░░░▄▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▄░░░░░▀░░")
                fmt.Println("     ░░░░░░░░█░░▄▀▀▀▀▀▀▀▀▀▀▀▀▀▄░░█░░░░░░░")
                fmt.Println("     ░░░░░░░░█░█░▀░░░░░▀░░▀░░░░█░█░░░░░░░")
                fmt.Println("     ░░▄▀▀█▄░█░█░░░░░░░░▄▀▀▄░▀░█░█▄▀▀▄░░░")
                fmt.Println("     ░░▀▄▄░▀██░█░░▀░░░░░█░░░▀▄▄█▄▀░░░█░░░")
                fmt.Println("     ░░░░▀█▄▄█░█▄░▀░░░▄▄▀░░░░░░░░░░░░▀▄░░")
                fmt.Println("     ░░░░░░░▀█░█░░░░▄░░█░░░▄█░░░▄░▄█░░█░░")
                fmt.Println("     ░░░░░░░░█░▀▄▀░░░░░█░██░▄░░▄░░▄░███░░")
                fmt.Println("     ░░░░░░░░█▄░░▀▀▀▀▀▀▀▀▄░░▀▀▀▀▀▀▀░▄▀░░░")
                fmt.Printf("     ░░░░░░░█░░▄█▀█▀▀█▀▀▀▀▀█▀▀█▀█▀▀█░░░░░   - You are elevator #1 (%v)\n", (*list_of_elevators)[0].Id)
                fmt.Println("     ░░▀░░░▀▀▀▀░░▀▀▀░░░░░░░▀▀▀░░▀▀░░░░▀░░░")
            
            } else if swap_counter >= 2 && swap_counter < 4 {
                fmt.Println("      ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ ")
                fmt.Println("     ░░░▀░░░░░▄▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▄░░░░░▀░░")
                fmt.Println("     ░░░░░░░░█░░▄▀▀▀▀▀▀▀▀▀▀▀▀▀▄░░█░░░░░░░")
                fmt.Println("     ░░░░░░░░█░█░▀░░░░░▀░░▀░░░░█░█░░░░░░░")
                fmt.Println("     ░░░░░░░░█░█░░░░░░░░▄▀▀▄░▀░█░█▄▀▀▄░░░")
                fmt.Println("     ░░░░░░░░█░█░░▀░░░░░█░░░▀▄▄█▄▀░░░█░░░")
                fmt.Println("     ░░░░░░░▄█░█▄░▀░░░▄▄▀░░░░░░░░░░░░▀▄░░")
                fmt.Println("     ░░▄▀▀▀▄ █░█░░░░▄░░█░░░▄█░░░▄░▄█░░█░░")
                fmt.Println("     ░░▀▄▄█▀▀█░▀▄▀░░░░░█░██░▄░░▄░░▄░███░░")
                fmt.Println("     ░░░░░░░░█▄░░▀▀▀▀▀▀▀▀▄░░▀▀▀▀▀▀▀░▄▀░░░")
                fmt.Printf("     ░░░░░░░█░░▄█▀█▀▀█▀▀▀▀▀█▀▀█▀█▀▀█░░░░░   - You are elevator #1 (%v)\n", (*list_of_elevators)[0].Id)
                fmt.Println("     ░░▀░░░▀▀▀▀░░▀▀▀░░░░░░░▀▀▀░░▀▀░░░░▀░░░")
            }

            fmt.Println("   ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░")
            fmt.Println("  ░░                                                                                               ░░")
            fmt.Printf("  ░░     Elevator name  \t\tStatus \t  Floor     Direction \tDoor open \tOrders \t   ░░\n")
            
            _mtx.Lock()
            for e:=0; e<len((*list_of_elevators)); e++ {
                fmt.Printf("  ░░  #%v %25v\t%v\t    %v\t    %6v\t%v\t\t  %d\t   ░░\n", 
                e+1,
                (*list_of_elevators)[e].Id, 
                (*list_of_elevators)[e].Status, 
                (*list_of_elevators)[e].Current_floor, 
                (*list_of_elevators)[e].Direction,
                (*list_of_elevators)[e].Door_open,
                len((*list_of_elevators)[e].Order_list))
            }
            _mtx.Unlock()
            
            fmt.Println("  ░░                                                                                               ░░")
            fmt.Println("   ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░")
            fmt.Printf("\n   - Network connections: %v", (*list_of_elevators)[0].Connections)
            fmt.Printf("\n   - Latest status msg: %s", new_status_msg)

            time.Sleep(90 * time.Millisecond)
            
            _mtx.Lock()
            if animate == true && (*list_of_elevators)[0].Status != "idle" && (*list_of_elevators)[0].Status != "error" { swap_counter++ }            
            _mtx.Unlock()

            if swap_counter > 8 { swap_counter = 0 }
        }
    }
}


/*********************************
* Print elevator states with json
**********************************/

func Print_status   (   list_of_elevators * []Elevator,
){

    time_asleep := 1 * time.Second

    for {
        fmt.Println("\n\n - ELEVATOR STATUS REPORT -")

        _mtx.Lock()
        data, _ := json.MarshalIndent((*list_of_elevators), "", "  ")
        _mtx.Unlock()

        fmt.Println(string(data))
        fmt.Println("")

        time.Sleep(time_asleep)
    }
}


/***********************************
* Function for sending dummy orders
************************************/

func Dummy_order    (   target_floor int,
                        order_buffer_ch chan Order,
){
    dummy_order := Order{}
    dummy_order.Id = 1
    dummy_order.Origin_e = "Dummy"
    dummy_order.Assigned_to = "Dummy"
    dummy_order.Button = 2
    dummy_order.Target_floor = target_floor

    order_buffer_ch <- dummy_order
}


/***************************************
* cost function, takes orders as input
* and gives order cost as output
****************************************/

func Order_cost(new_order Order, e Elevator) Cost{

    var cost_of_order Cost

    cost_of_order.Elevator_id = e.Id
    cost_of_order.Order_id = new_order.Id
    cost_of_order.Sum = 0

    // add cost by current status
    switch e.Status {
    case "idle":
        cost_of_order.Sum += 0
    case "busy":
        cost_of_order.Sum += 50
    case "error":
        cost_of_order.Sum += 100
    case "disconnected":
        cost_of_order.Sum += 1000
    default:
        cost_of_order.Sum += 0
    }  

    // add cost by distance to target
    cost_of_order.Sum += int(math.Abs(float64(new_order.Target_floor - e.Current_floor))) * 10 

    // add cost for direction - if target floor is in same direction as elevator is heading
    if (e.Direction == "up" && new_order.Target_floor < e.Current_floor) || (e.Direction == "down" && new_order.Target_floor > e.Current_floor){
        cost_of_order.Sum += 50
    }

    // add cost for each order in order list
    if e.Order_list != nil {
        for i := 0; i<len(e.Order_list); i++ {
            cost_of_order.Sum += 5
        }
    }

    return cost_of_order
}


/*********************************************
* Checks if order already exist in order list
**********************************************/

func Is_duplicate_order(elevators []Elevator, new_order Order) bool {

    is_duplicate_order := false

    _mtx.Lock()
    num_of_orders := len(elevators[0].Order_list)
    _mtx.Unlock()

    for i:=0; i<num_of_orders; i++ {     // run through order list to check if new order
        if elevators[0].Order_list[i].Button == new_order.Button && elevators[0].Order_list[i].Target_floor == new_order.Target_floor {
            is_duplicate_order = true
        }
    }

    return is_duplicate_order
}


/***********************************
* Add an order to a stack of orders
************************************/

func Push_order(item_to_add Order, stack []Order) []Order{
    return append(stack, item_to_add)
}


/********************************************
* Remove the first order from an order stack
*********************************************/

func Pop_order(stack []Order) []Order{    
    return stack[1:]
}


/********************************************
* Remove an indexed order from an order stack
*********************************************/

func Remove_order(index int, stack []Order) []Order{

    if len(stack) > 1{
        return append(stack[:index], stack[index+1:]...)
    } else {
        return nil
    }
}


/********************************************
* Remove an indexed order from an order stack
*********************************************/

func Remove_elevator(index int, stack []Elevator) []Elevator{

    if len(stack) > 1{
        return append(stack[:index], stack[index+1:]...)
    } else {
        return nil
    }
}


/*************************************
* Sync lights across all elevators by 
* cycling through all active orders
**************************************/

func Sync_lights    (   list_of_elevators * []Elevator,
){
    
    var buttons elevator_io.ButtonType

    for {

        // local stop button light      
        _mtx.Lock()
        if (*list_of_elevators)[0].Status == "error" || (*list_of_elevators)[0].Status == "Stop" {
            elevator_io.SetStopLamp(true)    // turn on stop light
        } else {            
            elevator_io.SetStopLamp(false)    // turn on stop light
        } 
        _mtx.Unlock()


        // sync hall lights
        for floors:=0; floors<4; floors++ { 

            for buttons=0; buttons<3; buttons++ {
                elevator_io.SetButtonLamp(buttons, floors, false)

                _mtx.Lock()
                num_of_elevators := len((*list_of_elevators))
                _mtx.Unlock()

                for e:=0; e<num_of_elevators; e++ { 

                    _mtx.Lock()
                    num_of_orders := len((*list_of_elevators)[e].Order_list)
                    _mtx.Unlock()

                    _mtx.Lock()
                    for o:=0; o<num_of_orders; o++ { 

                        if buttons == (*list_of_elevators)[e].Order_list[o].Button && buttons != 2 && floors == (*list_of_elevators)[e].Order_list[o].Target_floor {
                            elevator_io.SetButtonLamp(buttons, floors, true)

                        } else if buttons == (*list_of_elevators)[e].Order_list[o].Button && buttons == 2 && (*list_of_elevators)[e].Id == (*list_of_elevators)[0].Id && floors == (*list_of_elevators)[e].Order_list[o].Target_floor {
                            elevator_io.SetButtonLamp(buttons, floors, true)
                        }
                    }
                    _mtx.Unlock()
                }                
            }
        }
        time.Sleep(15 * time.Millisecond)
    }
}


/********************
* Open elevator door
*********************/

func Open_door  (   list_of_elevators * []Elevator,
){   
    time.Sleep(500 * time.Millisecond)      // wait for the passenger to exit
    
    elevator_io.SetDoorOpenLamp(true)            // turn on door lamp - door is opening
    
    _mtx.Lock()
    (*list_of_elevators)[0].Door_open = true           // set door status to open          
    _mtx.Unlock()

    time.Sleep(2000 * time.Millisecond)     // wait for the passenger to exit  
    
    elevator_io.SetDoorOpenLamp(false)           // turn off door lamp - door is closing
    
    _mtx.Lock()
    (*list_of_elevators)[0].Door_open = false          // set door status to closed        
    _mtx.Unlock()

    time.Sleep(1000 * time.Millisecond)     // wait a bit before next order  
}