import faker from 'faker'

const contactMock = () => {
    return {
        email: faker.internet.email(),
        extra: null,
        name: `${faker.name.firstName()} ${faker.name.lastName()}`,
        person: null,
        phone: faker.phone.phoneNumber()
    }
}

const contactTypeMock = () => faker.random.arrayElement(["Administrative", "Billing", "Legal", "Technical"])

export { contactMock, contactTypeMock }