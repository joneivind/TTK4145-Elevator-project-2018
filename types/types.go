package types

import  (
    "../elevator_io"
)

type Order struct {
    Id              int                     `json:"Order id"`          // order ID
    Origin_e        string                  `json:"From elevator"`     // From elevator
    Assigned_to     string                  `json:"Assigned to"`       // Assigned to
    Button          elevator_io.ButtonType  `json:"Button type"`       // button type
    Target_floor    int                     `json:"Target floor"`      // target floor
    Is_sent         bool                    `json:"Order sent"`        // order is sent to elevator
    Ack_received    bool                    `json:"ACK received"`      // ACK msg is received
    Error           int                     `json:"Error code"`        // error code
}

type Elevator struct {
    Id              string                  `json:"Id"`                // elevator IDs
    Status          string                  `json:"Status"`            // current elevator status
    Current_floor   int                     `json:"Current floor"`     // current floor
    Direction       string                  `json:"Direction"`         // current direction of elevator
    Door_open       bool                    `json:"Door open"`         // door status
    Connections     int                     `json:"Network connections"`// order list
    Order_list      []Order                 `json:"Order list"`        // order list
}

type Cost struct {
    Elevator_id     string
    Order_id        int
    Sum             int
}   