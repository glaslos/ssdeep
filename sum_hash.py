from peachpy import Argument, uint32_t
from peachpy.x86_64 import Function, GeneralPurposeRegister32, LOAD, MUL, RETURN, XOR, ADD
from peachpy.x86_64.registers import ax

f1 = Argument(uint32_t)
f2 = Argument(uint32_t)


with Function("sumHashASM", (f1, f2), uint32_t) as function:
    reg_f1 = GeneralPurposeRegister32()
    reg_f2 = GeneralPurposeRegister32()
    reg_f3 = GeneralPurposeRegister32()

    LOAD.ARGUMENT(reg_f1, f1)
    LOAD.ARGUMENT(reg_f2, f2)
    ADD(reg_f3, 0x01000193)

    ax = reg_f2

    MUL(reg_f3)
    XOR(ax, reg_f1)

    RETURN(ax)
