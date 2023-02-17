import { faker } from '@faker-js/faker';

import { AddressTypeHeaders } from '../constants';

const MockedGeographicAddress = () => [
  {
    address_line: faker.random.objectElement([[], [1, 2]]),
    address_type: faker.random.objectElement(Object.keys(AddressTypeHeaders)),
    building_name: '',
    building_number: '23',
    country: faker.address.country(),
    country_sub_division: 'MA',
    department: faker.commerce.department(),
    district_name: '',
    floor: '',
    post_box: '',
    post_code: faker.address.zipCode(),
    room: '',
    street_name: faker.address.streetName(),
    sub_department: faker.commerce.department(),
    town_location_name: '',
    town_name: faker.address.cityName(),
  },
];

export default MockedGeographicAddress;
