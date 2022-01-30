import faker from 'faker'

function generateOtherJuridictions(number = 2) {
    const juridictions = []

    for (let i = 1; i <= number; i++) {
        juridictions.push({
            country: faker.address.countryCode(),
            license_number: "",
            regulator_name: `${faker.name.firstName()} ${faker.name.lastName()}`
        })
    }

    return juridictions
}

export default generateOtherJuridictions