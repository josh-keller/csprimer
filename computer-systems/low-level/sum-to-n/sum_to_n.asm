global sum_to_n
section .text

sum_to_n:
		mov		rax, rdi		; move n to rax		
		inc		rdi
		mul		rdi
		sar		rax, 1
    ret
