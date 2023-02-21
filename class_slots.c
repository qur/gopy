// Code generated by gen_slots.py. DO NOT EDIT.

#include "utils.h"

#include "_cgo_export.h"

// ===============================================================

ClassContext *setSlots(PyTypeObject *type, uint64_t slotFlags) {
  ClassContext *ctxt = calloc(1, sizeof(ClassContext));
  ctxt->zero = NULL;

  type->tp_new = (newfunc)goClassNew;
  type->tp_dealloc = (destructor)goClassDealloc;

  {
    PyTypeObject *m = type;
    if (slotFlags & CLASS_HAS_TP_REPR) m->tp_repr = (reprfunc)goClassSlot_tp_repr;
    if (slotFlags & CLASS_HAS_TP_HASH) m->tp_hash = (hashfunc)goClassSlot_tp_hash;
    if (slotFlags & CLASS_HAS_TP_CALL) m->tp_call = (ternaryfunc)goClassSlot_tp_call;
    if (slotFlags & CLASS_HAS_TP_STR) m->tp_str = (reprfunc)goClassSlot_tp_str;
    if (slotFlags & CLASS_HAS_TP_GETATTRO) m->tp_getattro = (getattrofunc)goClassSlot_tp_getattro;
    if (slotFlags & CLASS_HAS_TP_SETATTRO) m->tp_setattro = (setattrofunc)goClassSlot_tp_setattro;
    if (slotFlags & CLASS_HAS_TP_RICHCOMPARE) m->tp_richcompare = (richcmpfunc)goClassSlot_tp_richcompare;
    if (slotFlags & CLASS_HAS_TP_ITER) m->tp_iter = (getiterfunc)goClassSlot_tp_iter;
    if (slotFlags & CLASS_HAS_TP_ITERNEXT) m->tp_iternext = (iternextfunc)goClassSlot_tp_iternext;
    if (slotFlags & CLASS_HAS_TP_DESCR_GET) m->tp_descr_get = (descrgetfunc)goClassSlot_tp_descr_get;
    if (slotFlags & CLASS_HAS_TP_DESCR_SET) m->tp_descr_set = (descrsetfunc)goClassSlot_tp_descr_set;
    if (slotFlags & CLASS_HAS_TP_INIT) m->tp_init = (initproc)goClassSlot_tp_init;
  }

  if (slotFlags & CLASS_HAS_AM) {
    PyAsyncMethods *m = &ctxt->am_meth;
    type->tp_as_async = m;
    if (slotFlags & CLASS_HAS_AM_AWAIT) m->am_await = (unaryfunc)goClassSlot_am_await;
    if (slotFlags & CLASS_HAS_AM_AITER) m->am_aiter = (unaryfunc)goClassSlot_am_aiter;
    if (slotFlags & CLASS_HAS_AM_ANEXT) m->am_anext = (unaryfunc)goClassSlot_am_anext;
    if (slotFlags & CLASS_HAS_AM_SEND) m->am_send = (sendfunc)goClassSlot_am_send;
  }

  if (slotFlags & CLASS_HAS_NB) {
    PyNumberMethods *m = &ctxt->nb_meth;
    type->tp_as_number = m;
    if (slotFlags & CLASS_HAS_NB_ADD) m->nb_add = (binaryfunc)goClassSlot_nb_add;
    if (slotFlags & CLASS_HAS_NB_INPLACE_ADD) m->nb_inplace_add = (binaryfunc)goClassSlot_nb_inplace_add;
    if (slotFlags & CLASS_HAS_NB_SUBTRACT) m->nb_subtract = (binaryfunc)goClassSlot_nb_subtract;
    if (slotFlags & CLASS_HAS_NB_INPLACE_SUBTRACT) m->nb_inplace_subtract = (binaryfunc)goClassSlot_nb_inplace_subtract;
    if (slotFlags & CLASS_HAS_NB_MULTIPLY) m->nb_multiply = (binaryfunc)goClassSlot_nb_multiply;
    if (slotFlags & CLASS_HAS_NB_INPLACE_MULTIPLY) m->nb_inplace_multiply = (binaryfunc)goClassSlot_nb_inplace_multiply;
    if (slotFlags & CLASS_HAS_NB_REMAINDER) m->nb_remainder = (binaryfunc)goClassSlot_nb_remainder;
    if (slotFlags & CLASS_HAS_NB_INPLACE_REMAINDER) m->nb_inplace_remainder = (binaryfunc)goClassSlot_nb_inplace_remainder;
    if (slotFlags & CLASS_HAS_NB_DIVMOD) m->nb_divmod = (binaryfunc)goClassSlot_nb_divmod;
    if (slotFlags & CLASS_HAS_NB_POWER) m->nb_power = (ternaryfunc)goClassSlot_nb_power;
    if (slotFlags & CLASS_HAS_NB_INPLACE_POWER) m->nb_inplace_power = (ternaryfunc)goClassSlot_nb_inplace_power;
    if (slotFlags & CLASS_HAS_NB_NEGATIVE) m->nb_negative = (unaryfunc)goClassSlot_nb_negative;
    if (slotFlags & CLASS_HAS_NB_POSITIVE) m->nb_positive = (unaryfunc)goClassSlot_nb_positive;
    if (slotFlags & CLASS_HAS_NB_ABSOLUTE) m->nb_absolute = (unaryfunc)goClassSlot_nb_absolute;
    if (slotFlags & CLASS_HAS_NB_BOOL) m->nb_bool = (inquiry)goClassSlot_nb_bool;
    if (slotFlags & CLASS_HAS_NB_INVERT) m->nb_invert = (unaryfunc)goClassSlot_nb_invert;
    if (slotFlags & CLASS_HAS_NB_LSHIFT) m->nb_lshift = (binaryfunc)goClassSlot_nb_lshift;
    if (slotFlags & CLASS_HAS_NB_INPLACE_LSHIFT) m->nb_inplace_lshift = (binaryfunc)goClassSlot_nb_inplace_lshift;
    if (slotFlags & CLASS_HAS_NB_RSHIFT) m->nb_rshift = (binaryfunc)goClassSlot_nb_rshift;
    if (slotFlags & CLASS_HAS_NB_INPLACE_RSHIFT) m->nb_inplace_rshift = (binaryfunc)goClassSlot_nb_inplace_rshift;
    if (slotFlags & CLASS_HAS_NB_AND) m->nb_and = (binaryfunc)goClassSlot_nb_and;
    if (slotFlags & CLASS_HAS_NB_INPLACE_AND) m->nb_inplace_and = (binaryfunc)goClassSlot_nb_inplace_and;
    if (slotFlags & CLASS_HAS_NB_XOR) m->nb_xor = (binaryfunc)goClassSlot_nb_xor;
    if (slotFlags & CLASS_HAS_NB_INPLACE_XOR) m->nb_inplace_xor = (binaryfunc)goClassSlot_nb_inplace_xor;
    if (slotFlags & CLASS_HAS_NB_OR) m->nb_or = (binaryfunc)goClassSlot_nb_or;
    if (slotFlags & CLASS_HAS_NB_INPLACE_OR) m->nb_inplace_or = (binaryfunc)goClassSlot_nb_inplace_or;
    if (slotFlags & CLASS_HAS_NB_INT) m->nb_int = (unaryfunc)goClassSlot_nb_int;
    if (slotFlags & CLASS_HAS_NB_FLOAT) m->nb_float = (unaryfunc)goClassSlot_nb_float;
    if (slotFlags & CLASS_HAS_NB_FLOOR_DIVIDE) m->nb_floor_divide = (binaryfunc)goClassSlot_nb_floor_divide;
    if (slotFlags & CLASS_HAS_NB_INPLACE_FLOOR_DIVIDE) m->nb_inplace_floor_divide = (binaryfunc)goClassSlot_nb_inplace_floor_divide;
    if (slotFlags & CLASS_HAS_NB_TRUE_DIVIDE) m->nb_true_divide = (binaryfunc)goClassSlot_nb_true_divide;
    if (slotFlags & CLASS_HAS_NB_INPLACE_TRUE_DIVIDE) m->nb_inplace_true_divide = (binaryfunc)goClassSlot_nb_inplace_true_divide;
    if (slotFlags & CLASS_HAS_NB_INDEX) m->nb_index = (unaryfunc)goClassSlot_nb_index;
    if (slotFlags & CLASS_HAS_NB_MATRIX_MULTIPLY) m->nb_matrix_multiply = (binaryfunc)goClassSlot_nb_matrix_multiply;
    if (slotFlags & CLASS_HAS_NB_INPLACE_MATRIX_MULTIPLY) m->nb_inplace_matrix_multiply = (binaryfunc)goClassSlot_nb_inplace_matrix_multiply;
  }

  if (slotFlags & CLASS_HAS_MP) {
    PyMappingMethods *m = &ctxt->mp_meth;
    type->tp_as_mapping = m;
    if (slotFlags & CLASS_HAS_MP_LENGTH) m->mp_length = (lenfunc)goClassSlot_mp_length;
    if (slotFlags & CLASS_HAS_MP_SUBSCRIPT) m->mp_subscript = (binaryfunc)goClassSlot_mp_subscript;
    if (slotFlags & CLASS_HAS_MP_ASS_SUBSCRIPT) m->mp_ass_subscript = (objobjargproc)goClassSlot_mp_ass_subscript;
  }

  if (slotFlags & CLASS_HAS_SQ) {
    PySequenceMethods *m = &ctxt->sq_meth;
    type->tp_as_sequence = m;
    if (slotFlags & CLASS_HAS_SQ_LENGTH) m->sq_length = (lenfunc)goClassSlot_sq_length;
    if (slotFlags & CLASS_HAS_SQ_CONCAT) m->sq_concat = (binaryfunc)goClassSlot_sq_concat;
    if (slotFlags & CLASS_HAS_SQ_REPEAT) m->sq_repeat = (ssizeargfunc)goClassSlot_sq_repeat;
    if (slotFlags & CLASS_HAS_SQ_ITEM) m->sq_item = (ssizeargfunc)goClassSlot_sq_item;
    if (slotFlags & CLASS_HAS_SQ_ASS_ITEM) m->sq_ass_item = (ssizeobjargproc)goClassSlot_sq_ass_item;
    if (slotFlags & CLASS_HAS_SQ_CONTAINS) m->sq_contains = (objobjproc)goClassSlot_sq_contains;
    if (slotFlags & CLASS_HAS_SQ_INPLACE_CONCAT) m->sq_inplace_concat = (binaryfunc)goClassSlot_sq_inplace_concat;
    if (slotFlags & CLASS_HAS_SQ_INPLACE_REPEAT) m->sq_inplace_repeat = (ssizeargfunc)goClassSlot_sq_inplace_repeat;
  }

  if (slotFlags & CLASS_HAS_BF) {
    PyBufferProcs *m = &ctxt->bf_meth;
    type->tp_as_buffer = m;
    if (slotFlags & CLASS_HAS_BF_GETBUFFER) m->bf_getbuffer = (getbufferproc)goClassSlot_bf_getbuffer;
    if (slotFlags & CLASS_HAS_BF_RELEASEBUFFER) m->bf_releasebuffer = (releasebufferproc)goClassSlot_bf_releasebuffer;
  }

  return ctxt;
}

// ===============================================================
