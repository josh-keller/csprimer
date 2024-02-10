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


def average_age(ages):
    total = 0
    for a in ages:
        total += a
    return total / len(ages)


def average_payment_amount(amounts):
    amount = 0
    for a in amounts:
        amount += a
    return (amount / 100) / len(amounts)


def stddev_payment_amount(amounts):
    mean = average_payment_amount(amounts)
    squared_diffs = 0
    for a in amounts:
        amount = a / 100
        diff = amount - mean
        squared_diffs += diff * diff
    return math.sqrt(squared_diffs / len(amounts))


def load_data():
    users = {}
    amounts = []
    ages = []

    with open('users.csv') as f:
        for line in csv.reader(f):
            uid, name, age, address_line, zip_code = line
            addr = Address(address_line, zip_code)
            age = int(age)
            users[int(uid)] = User(int(uid), name, age, addr, [])
            ages.append(age)

    with open('payments.csv') as f:
        for line in csv.reader(f):
            amount, timestamp, uid = line
            amount = int(amount)
            dollarAmount = DollarAmount(dollars=amount//100, cents=amount % 100)
            payment = Payment(
                dollarAmount,
                time=datetime.datetime.fromisoformat(timestamp))
            users[int(uid)].payments.append(payment)
            amounts.append(amount)
    return users, ages, amounts


if __name__ == '__main__':
    t = time.perf_counter()
    users, ages, amounts = load_data()
    print(f'Data loading: {time.perf_counter() - t:.3f}s')
    t = time.perf_counter()
    assert abs(average_age(ages) - 59.626) < 0.01
    assert abs(stddev_payment_amount(amounts) - 288684.849) < 0.01
    assert abs(average_payment_amount(amounts) - 499850.559) < 0.01
    print(f'Computation {time.perf_counter() - t:.3f}s')
