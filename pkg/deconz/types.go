package deconz

type Light struct {
	Etag             string      `json:"etag,omitempty"`
	Hascolor         bool        `json:"hascolor,omitempty"`
	Lastannounced    interface{} `json:"lastannounced,omitempty"`
	Lastseen         string      `json:"lastseen,omitempty"`
	Manufacturername string      `json:"manufacturername,omitempty"`
	Modelid          string      `json:"modelid,omitempty"`
	Name             string      `json:"name"`
	State            State       `json:"state,omitempty"`
	Swversion        string      `json:"swversion,omitempty"`
	Type             string      `json:"type,omitempty"`
	Uniqueid         string      `json:"uniqueid,omitempty"`
}

type State struct {
	Alert     string `json:"alert,omitempty"`
	Bri       int    `json:"bri,omitempty"`
	On        bool   `json:"on,omitempty"`
	Reachable bool   `json:"reachable,omitempty"`
	Colormode string `json:"colormode,omitempty"`
	Ct        int    `json:"ct,omitempty"`
}
