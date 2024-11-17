import csv
import datetime
import math
import time


class Address(object):
    def __init__(self, address_line, zipcode):
        self.address_line = address_line
        self.zipcode = zipcode


class DollarAmount(object):
    def __init__(self, dollars, cents):
        self.dollars = dollars
        self.cents = cents


class Payment(object):
    def __init__(self, dollar_amount, time):
        self.amount = dollar_amount
        self.time = time


class User(object):
    def __init__(self, user_id, name, age, address, payments):
        self.user_id = user_id
        self.name = name
        self.age = age
        self.address = address
        self.payments = payments


def average_age(users):
    total = 0
    for u in users.values():
        total += u.age
    return total / len(users)


def average_payment_amount(users):
    amount = 0
    count = 0
    for u in users.values():
        count += len(u.payments)
        for p in u.payments:
            amount += float(p.amount.dollars) + float(p.amount.cents) / 100
    return amount / count


def stddev_payment_amount(users):
    mean = average_payment_amount(users)
    squared_diffs = 0
    count = 0
    for u in users.values():
        count += len(u.payments)
        for p in u.payments:
            amount = float(p.amount.dollars) + float(p.amount.cents) / 100
            diff = amount - mean
            squared_diffs += diff * diff
    return math.sqrt(squared_diffs / count)


def load_data():
    users = {}
    with open('users.csv') as f:
        for line in csv.reader(f):
            uid, name, age, address_line, zip_code = line
            addr = Address(address_line, zip_code)
            users[int(uid)] = User(int(uid), name, int(age), addr, [])
    with open('payments.csv') as f:
        for line in csv.reader(f):
            amount, timestamp, uid = line
            payment = Payment(
                DollarAmount(dollars=int(amount)//100, cents=int(amount) % 100),
                time=datetime.datetime.fromisoformat(timestamp))
            users[int(uid)].payments.append(payment)
    return users


if __name__ == '__main__':
    t = time.perf_counter()
    users = load_data()
    print(f'Data loading: {time.perf_counter() - t:.3f}s')
    t = time.perf_counter()
    assert abs(average_age(users) - 59.626) < 0.01
    assert abs(stddev_payment_amount(users) - 288684.849) < 0.01
    assert abs(average_payment_amount(users) - 499850.559) < 0.01
    print(f'Computation {time.perf_counter() - t:.3f}s')
