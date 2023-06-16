package main

/***************************
    @author: tiansheng.ren
    @date: 6/3/23
    @desc:

***************************/

type condition struct {
	field string
	op    string
	value interface{}
}

func NewCondition() *condition {
	return &condition{}
}

func (c *condition) Set(k, op string, v interface{}) {
	c.field, c.op, c.value = k, op, v
}
