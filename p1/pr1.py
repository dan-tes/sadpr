from datetime import datetime, timedelta
from abc import ABC, ABCMeta
from collections import defaultdict

class EntityMeta(ABCMeta):
    def __repr__(cls):
        return f"<Класс {cls.__name__}>"


class Entity(ABC, metaclass=EntityMeta):
    """Базовый класс для всех сущностей"""

    def __init__(self):
        pass

    def __repr__(self):
        return f"<Класс {self.__name__}>"


class Employee(Entity):
    def __init__(self, name: str, salary: float):
        super().__init__()
        self.name = name
        self.salary = salary

    def __repr__(self):
        return f"Сотрудник: {self.name}, зарплата: {self.salary}"


class CallCenterOperator(Employee):
    def __repr__(self):
        return f"Оператор колл-центра (звонки): {self.name}"


class MessengerOperator(Employee):
    def __repr__(self):
        return f"Оператор колл-центра (мессенджер): {self.name}"


class Admin(Employee):
    def __init__(self, name: str, salary: float, responsible_for: Entity):
        super().__init__(name, salary)
        self.responsible_for = responsible_for

    def __repr__(self):
        return f"Администратор: {self.name}, отвечает за: {self.responsible_for}"


# --- Клиенты ---

class Client(Entity):
    def __init__(self, iin: str = None):
        super().__init__()
        self.iin = iin

    def __repr__(self):
        return f"Клиент: ИИН: {self.iin}"


class Company(Client):
    def __init__(self, name: str, iin: str = None):
        super().__init__(iin)
        self.name = name

    def __repr__(self):
        return f'Компания {self.name}, ИИН: {self.iin}'


class Individual(Client):
    def __init__(self, fio: str, iin: str = None):
        super().__init__(iin)
        self.fio = fio

    def __repr__(self):
        return f'Индивидуальный клиент {self.fio}, ИИН: {self.iin}'


class Accountant(Employee):
    def __repr__(self):
        return f"Бухгалтер: {self.name}"


# --- Продукты ---

class Product(Entity):
    def __init__(self, price: float, article: int):
        super().__init__()
        self.price = price
        self.article = article

    def __repr__(self):
        return f"Продукт: артикул {self.article}, цена {self.price} руб."


# --- Каналы связи ---

class CommunicationChannel(Entity):
    def __init__(self, sender: Client, time_received: datetime):
        super().__init__()
        self.sender = sender
        self.time_received = time_received


class Message(CommunicationChannel):
    def __init__(self, sender: Client, text: str, time_received: datetime):
        super().__init__(sender, time_received)
        self.text = text

    def __repr__(self):
        return f"Сообщение от {self.sender}: {self.text}"


class Call(CommunicationChannel):
    def __init__(self, sender: Client, duration: int, transcript: str, time_received: datetime):
        super().__init__(sender, time_received)
        self.duration = duration
        self.transcript = transcript

    def __repr__(self):
        return f"Звонок от {self.sender}, длительность {self.duration} сек., расшифровка: {self.transcript}"


# --- Заказы ---

class Order(Entity):
    def __init__(self, time: datetime, client: Client, products: list[Product], summ: int):
        super().__init__()
        self.time = time
        self.client = client
        self.products = products
        self.summ = summ

    def __repr__(self):
        return f"Заказ от {self.time}, клиент: {self.client}, сумма: {self.summ} руб."


class RegularOrder(Order):
    def __init__(self, time: datetime, client: Client, products: list[Product], summ: int, del_time: timedelta):
        super().__init__(time, client, products, summ)
        self.del_time = del_time

    def __repr__(self):
        return f"Регулярный заказ: {self.time}, клиент: {self.client}, сумма: {self.summ}, период {self.del_time}"


class OneTimeOrder(Order):
    def __init__(self, time: datetime, client: Client, products: list[Product], summ: int):
        super().__init__(time, client, products, summ)

    def __repr__(self):
        return f"Разовый заказ: {self.time}, клиент: {self.client}, сумма: {self.summ}"







def create_database():
    db = defaultdict(list)  # ключ = класс, значение = список экземпляров

    # --- Сотрудники ---
    emp1 = CallCenterOperator("Иван", 40000)
    emp2 = CallCenterOperator("Алексей", 38000)
    emp3 = MessengerOperator("Мария", 42000)
    emp4 = MessengerOperator("Светлана", 41000)
    admin1 = Admin("Олег", 60000, Company)
    admin2 = Admin("Никита", 55000, Individual)
    acc1 = Accountant("Анна", 50000)
    acc2 = Accountant("Екатерина", 52000)

    # --- Клиенты ---
    client1 = Company("TechCorp", "123456789012")
    client2 = Company("SoftInc", "234567890123")
    client3 = Individual("Петров Петр Петрович", "987654321098")
    client4 = Individual("Иванова Ирина Ивановна", "876543210987")

    # --- Продукты ---
    product1 = Product(1500, 101)
    product2 = Product(2500, 202)
    product3 = Product(500, 303)
    product4 = Product(750, 404)

    # --- Сообщения и звонки ---
    msg1 = Message(client1, "Хочу заказать продукцию", datetime.now())
    msg3 = Message(client1, "ОЧЕНЬОЧЕНЬОЧЕНЬОЧЕНЬОЧЕНЬОЧЕНЬ Хочу заказать продукцию ", datetime.now())

    msg2 = Message(client2, "Есть ли скидки?", datetime.now())
    call1 = Call(client3, 120, "Обсуждение условий доставки", datetime.now())
    call2 = Call(client4, 90, "Вопрос по оплате", datetime.now())

    # --- Заказы ---
    order1 = OneTimeOrder(datetime.now(), client1, [product1, product2], 4000)
    order2 = OneTimeOrder(datetime.now(), client2, [product3], 500)
    order3 = RegularOrder(datetime.now(), client3, [product2, product4], 3250, timedelta(days=30))
    order4 = RegularOrder(datetime.now(), client4, [product1], 1500, timedelta(days=7))

    # --- Добавляем все в базу ---
    for obj in [emp1, emp2, emp3, emp4, admin1, admin2, acc1, acc2]:
        db[type(obj)].append(obj)

    for obj in [client1, client2, client3, client4]:
        db[type(obj)].append(obj)

    for obj in [product1, product2, product3, product4]:
        db[type(obj)].append(obj)

    for obj in [msg1, msg2, call1, call2,msg3]:
        db[type(obj)].append(obj)

    for obj in [order1, order2, order3, order4]:
        db[type(obj)].append(obj)

    return db


def find(cls, attr, question, val):
    print('Итоговый запрос', cls, attr, question, val)
    arr = database[cls]
    final_list = []
    for obj in arr:
        func_list = {'==': lambda a, b: a == b, '!=': lambda a, b: str(b) not in str(a), ">=": lambda a, b: a >= b,
                     "<=": lambda a, b: a <= b, '>': lambda a, b: a > b, '<': lambda a, b: a < b}
        spec_func = {'==': lambda a, b: a in b, '!=': lambda a, b: a not in b}
        attrib = getattr(obj, attr)
        if isinstance(attrib, EntityMeta) and spec_func[question](val, str(attrib)):
            final_list.append(obj)

        elif func_list[question](attrib, val):
            final_list.append(obj)
    print(final_list)


if __name__ == "__main__":
    # --- Использование ---
    database = create_database()
    find(Company, 'name', '==', 'TechCorp')
    print('1:1')
    find(Admin, 'responsible_for', '!=', 'Company')
    print('long')
    find(Message, 'text', '!=', 'Хочу акцию')
    print('1:9010192938')
    while True:
        user_input = input("Введите название класса: ")

        # достанем класс по имени
        cls = globals().get(user_input)

        if cls is None:
            print("Такого класса нет!")
            continue
        else:
            print(f"Вы выбрали класс: {cls}")


            def create_instance_with_defaults(cls):
                try:
                    if issubclass(cls, Company):
                        return cls(name="TestCorp", iin="123456789012")
                    elif issubclass(cls, Individual):
                        return cls(fio="Иванов Иван", iin="987654321098")
                    elif issubclass(cls, Employee):
                        if issubclass(cls, Admin):
                            return cls(name="Test", salary=0, responsible_for=Employee)
                        return cls(name="Test", salary=0)
                    elif issubclass(cls, Product):
                        return cls(price=0, article=0)
                    elif issubclass(cls, Message):
                        dummy_client = Client("000000000000")
                        return cls(sender=dummy_client, text="Тест", time_received=datetime.now())
                    elif issubclass(cls, Call):
                        dummy_client = Client("000000000000")
                        return cls(sender=dummy_client, duration=0, transcript="Тест", time_received=datetime.now())
                    elif issubclass(cls, Order):
                        dummy_client = Client("000000000000")
                        dummy_product = Product(0, 0)
                        if issubclass(cls, RegularOrder):
                            return cls(datetime.now(), dummy_client, [dummy_product], 0, timedelta(days=30))
                        else:
                            return cls(datetime.now(), dummy_client, [dummy_product], 0)
                    elif issubclass(cls, Client):
                        return cls(iin="000000000000")
                    else:
                        return cls()  # конструктор без аргументов
                except TypeError as e:
                    print(f"Не удалось создать экземпляр: {e}")
                    return None


            instance = create_instance_with_defaults(cls)
            if instance:
                print(f"Экземпляр класса {cls.__name__}: {instance}")
                attrs = '\n'.join([f'{i} тип {type(i)}' for i in instance.__dict__.keys()])
                print(f"Атрибуты экземпляра: \n{attrs}")
            else:
                print("Не удалось создать экземпляр для получения атрибутов.")
        attr = input('Введите поле ')
        if attr not in [str(i) for i in instance.__dict__.keys()]:
            print('Такого поля нет')
            continue
        print('Выбрано поле ', attr)
        print('Как сравнивать? (== != for int, date >= <= > <) ')

        question = input()
        if not (question in ['==', '!='] or (question in ['>=', '<=' ,'>' ,'<'] and not type(attr) is str)):
            print('Нельзя')
            continue
        val = input('Введите значение ')

        find(cls, attr, question, val)
