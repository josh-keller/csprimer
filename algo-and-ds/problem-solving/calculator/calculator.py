def evaluate(stack):
    b = stack.pop()
    op = stack.pop()
    a = stack.pop()

    if op == "+":
        return a + b

    if op == "-":
        return a - b

    if op == "/":
        return a / b

    if op == "*":
        return a * b

    return None


def push_num(num, stack):
    stack.append(int("".join(num)))
    num.clear()


def stack_eval(exp):
    stack = []
    num = []

    for char in exp:
        # If char is a digit, append to num and move on
        if ord(char) in range(48, 58):
            num.append(char)
            continue

        # If char is not a number and we have something in the num accumulator, handle it
        if num:
            push_num(num, stack)

        if char in ["(", " "]:             # Do nothing for open paren or space
            continue
        elif char == ")":                  # Close paren means we have 3 values ready to eval in stack
            stack.append(evaluate(stack))
        elif char in ['+', '-', '/', '*']: # Op just needs to be pushed to stack
            stack.append(char)
        else:
            raise("Parsing error")

    return stack.pop()


assert(stack_eval("(11 + 2)") == 13)
assert(stack_eval("(1 - (2 - 3))") == 2)
assert(stack_eval("(1 + (2-3))") == 0)
assert(stack_eval("(1 + ((3*4) - 3))") == 10)
assert(stack_eval("((8/4) - ((3*41) - 32))") == -89)
print('ok')
