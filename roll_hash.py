from peachpy import Argument, uint32_t, ptr
from peachpy.x86_64 import Function, GeneralPurposeRegister64, GeneralPurposeRegister32, LOAD, MOVZX, ADD, RETURN, word, SUB, MUL
from peachpy.x86_64.registers import ax

f1 = Argument(ptr())
f2 = Argument(uint32_t)

with Function("rollHashASM", (f1, f2), uint32_t) as function:
    """
    rollingState.window []byte: 0-24 (size 24, align 8)
    rollingState.h1 uint32: 24-28 (size 4, align 4)
    rollingState.h2 uint32: 28-32 (size 4, align 4)
    rollingState.h3 uint32: 32-36 (size 4, align 4)
    rollingState.n uint32: 36-40 (size 4, align 4)
    """

    reg_rolling_state = GeneralPurposeRegister64()
    reg_byte = GeneralPurposeRegister32()
    reg_h1 = GeneralPurposeRegister32()
    reg_h2 = GeneralPurposeRegister32()
    reg_h3 = GeneralPurposeRegister32()
    reg_n = GeneralPurposeRegister32()

    LOAD.ARGUMENT(reg_rolling_state, f1)
    LOAD.ARGUMENT(reg_byte, f2)

    MOVZX(reg_h1, word[reg_rolling_state+24])
    MOVZX(reg_h2, word[reg_rolling_state+28])
    MOVZX(reg_h3, word[reg_rolling_state+32])
    MOVZX(reg_n, word[reg_rolling_state+36])

    # rs.h2 -= rs.h1
    SUB(reg_h2, reg_h1)

	# rs.h2 += rollingWindow * uint32(c)
    ax = GeneralPurposeRegister32()
    ADD(ax, 7)
    MUL(reg_byte)
    ADD(reg_h2, ax)

	# rs.h1 += uint32(c)
    ADD(reg_h1, reg_byte)

	# rs.h1 -= uint32(rs.window[rs.n])
    # rs.window[rs.n] = c
    # TODO: SUB(reg_h1, )

	# rs.n++
    ADD(reg_n, 1)

    #if rs.n == rollingWindow {
	#	rs.n = 0
	#}
    #rs.h3 = rs.h3 << 5
	#rs.h3 ^= uint32(c)
    # TODO:

    # return rs.h1 + rs.h2 + rs.h3
    ADD(reg_h1, reg_h2)
    ADD(reg_h1, reg_h3)

    RETURN(reg_h1)
