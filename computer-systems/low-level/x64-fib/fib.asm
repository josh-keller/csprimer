section .text
global fib
fib:
  mov			rax, edi
	cmp 		edi, 1          ; if n is <= 1
	jle		 	.return      		; return it

	dec    	edi             ; derement n
	push   	edi						  ; push to stack
	call   	fib							; call fib with n-1 (in edi)
	pop		 	edi							; fib(n-1) is now in rax, pop n-1 back into rdi

	push   	rax							; push rax to stack for later
	dec    	edi							; decrement n again (n-2)
	call   	fib							; call fib with n-2 (in edi)
	pop		 	edi							; pop it back 

	add		 	rax, edi
.return:
	ret
