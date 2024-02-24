package main

import (
	"testing"
)

var testMessages []message = []message{
	{
		tail_number:  "N20904",
		engine_count: 2,
		engine_name:  "GEnx-1B",
		latitude:     39.11593389482025,
		longitude:    -67.32425341289998,
		altitude:     36895.5,
		temperature:  -53.2,
	},
	{
		tail_number:  "N20906",
		engine_count: 2,
		engine_name:  "GEnx-1B",
		latitude:     83.31593389482026,
		longitude:    -7.12425341290001,
		altitude:     16895.5,
		temperature:  -13.2,
	},
	{
		tail_number:  "N20907",
		engine_count: 4,
		engine_name:  "GEnx-1C",
		latitude:     -3.31593389482026,
		longitude:    17.12425341290001,
		altitude:     7032.5,
		temperature:  0.2,
	},
	{
		tail_number:  "",
		engine_count: 0,
		engine_name:  "",
		latitude:     0,
		longitude:    0,
		altitude:     0,
		temperature:  0,
	},
}

func TestBinaryMarshaling(t *testing.T) {
	for _, msg := range testMessages {
		bin, err := msg.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}
		//fmt.Println(hex.Dump(out))
		mp := new(message)
		err = mp.UnmarshalBinary(bin)
		if err != nil {
			t.Fatal(err)
		}
		newmsg := *mp
		//fmt.Printf("%#v\n", newmsg)
		if msg != newmsg {
			t.Errorf("bad message; want: `%#v`, got: `%#v`", msg, newmsg)
		}
	}
}
