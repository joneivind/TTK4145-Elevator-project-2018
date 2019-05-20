package main

import  (
    "./order"
    ."./types"
    "./elevator_io"
    "./network/bcast"
    "./network/peers"
    "./network/nethandler"
    "./misc"
    "sync"
)

var _mtx sync.Mutex

func main(){

    // create list of elevators
    var elevators []Elevator  


    network_name        := "nyan"       // network tag
	numFloors           := 4            // number of floors
	elevator_bcast_rate := 100          // ms delay between transmitting elevator state
	max_channel_input   := 10           // limit spam on channels

    elevator_adress     := "localhost:15657"
	port_peers          := 15645
	port_bcast_orders   := 16569
	port_bcast_elevator := 16500 


    // io channels
    drv_buttons     := make(chan elevator_io.ButtonEvent)
    drv_floors      := make(chan int)
    drv_obstr       := make(chan bool)
    drv_stop        := make(chan bool)

    // order channels
    orders_ch       := make(chan Order, max_channel_input)
    order_buffer_ch := make(chan Order, max_channel_input)
    order_ack_ch    := make(chan int)
    order_error_ch  := make(chan int)
    failed_order_ch := make(chan Order)
    order_tx_ch     := make(chan Order, max_channel_input)
    order_rx_ch     := make(chan Order, max_channel_input)

    // elevator state channels
    elevator_tx_ch  := make(chan Elevator)
    elevator_rx_ch  := make(chan Elevator, max_channel_input)

    // network peer update channels
    peer_update_ch  := make(chan peers.PeerUpdate)
    peer_tx_enable  := make(chan bool)

    status_msg      := make(chan string, max_channel_input)



    // initialize elevator
    elevator_io.Init(elevator_adress, numFloors)
    elevator_io.Reset_lights()
  	

    // create a new elevator
    nethandler.Create_new_local_elevator	(	network_name,
                                    			&elevators)  

    // start io routines
    go elevator_io.PollButtons(drv_buttons)
    go elevator_io.PollFloorSensor(drv_floors)
    go elevator_io.PollObstructionSwitch(drv_obstr)
    go elevator_io.PollStopButton(drv_stop)    

    // sync lights across all elevators
    go misc.Sync_lights		(	&elevators)

    // start elevator order routines
    go order.Order_assigner	(   drv_buttons,
                                order_buffer_ch, 
                                order_tx_ch, 
                                failed_order_ch,
                                status_msg,
                                &elevators) 

    go order.Order_handler	(   order_buffer_ch, 
                                orders_ch, 
                                order_ack_ch, 
                                order_error_ch, 
                                order_tx_ch, 
                                failed_order_ch,
                                status_msg,
                                &elevators)
        
    go order.Order_executer	(   orders_ch, 
                                order_ack_ch, 
                                drv_floors, 
                                drv_stop, 
                                order_error_ch,
                                status_msg,
                                &elevators)

    go order.Add_cab_orders_from_file	(	order_rx_ch, 
		                                    numFloors,
		                                    status_msg,
		                                    &elevators)

    // network updates and receiver handling
    go nethandler.Peer_update_handler	(	peer_update_ch, 
			                                failed_order_ch,
			                                network_name,
			                                status_msg,
			                                &elevators)

    go nethandler.Network_order_rx     (   order_buffer_ch, 
			                                order_rx_ch,
			                                status_msg,
			                                &elevators)

    go nethandler.Elevator_state_tx    (   elevator_bcast_rate, 
			                                elevator_tx_ch,
			                                &elevators)

    go nethandler.Elevator_state_rx    (   elevator_rx_ch,
                                			&elevators)

    go peers.Transmitter    (   port_peers, 
                                elevators[0].Id, 
                                peer_tx_enable)

    go peers.Receiver       (   port_peers, 
                                peer_update_ch)

    go bcast.Transmitter    (   port_bcast_orders, 
                                order_tx_ch)

    go bcast.Receiver       (   port_bcast_orders, 
                                order_rx_ch)

    go bcast.Transmitter    (   port_bcast_elevator, 
                                elevator_tx_ch)

    go bcast.Receiver       (   port_bcast_elevator, 
                                elevator_rx_ch)

    go misc.Print_state	(	status_msg, 
    						drv_obstr,
                        	&elevators)

    // loop program
    for { select {} }
}
