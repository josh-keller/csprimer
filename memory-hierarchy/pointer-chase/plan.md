# Pointer Chase

Okay, so the data is currently structured like this:

Users: [
  {
    user_id,
    name: string,
    age: number,
    address: Address{
        address_line: string,
        zipcode: string? number?
    },
    payments: []Payment{
        amount: number,
        time: datetime
    }
  }
]

Let's start with what kind of speed up we can get, just by splitting the payments out.
I don't think there's a good reason to have the payments nested, other than quick lookup.
Two options for payments:
1. A list (slow lookup)
2. A dict (quicker lookup, but what key are we looking up on?)

