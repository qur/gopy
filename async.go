package py

// #include "utils.h"
import "C"

type SendResult int

const (
	SendReturn SendResult = C.PYGEN_RETURN
	SendError  SendResult = C.PYGEN_ERROR
	SendNext   SendResult = C.PYGEN_NEXT
)
