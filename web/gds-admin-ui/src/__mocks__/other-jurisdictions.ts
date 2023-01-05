import {faker} from '@faker-js/faker';

function generateOtherJurisdictions(number = 2) {
  const jurisdictions: any = [];

  for (let i = 1; i <= number; i++) {
    jurisdictions.push({
      country: 'US',
      license_number: '',
      regulator_name: `${faker.name.firstName()} ${faker.name.lastName()}`,
    });
  }

  return jurisdictions;
}

export default generateOtherJurisdictions;
