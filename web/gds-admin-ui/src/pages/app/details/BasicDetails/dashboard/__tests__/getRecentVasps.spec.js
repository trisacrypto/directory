import { getRecentVasps } from "../Tasks"
import faker from "faker"

const FAKE_SYSTEM_TIMESTAMP = 1638305340000; // 2021-11-30T20:49:00.000Z

describe('getRecentVasps', () => {

    beforeAll(() => {
        jest.useFakeTimers('modern').setSystemTime(FAKE_SYSTEM_TIMESTAMP);
    });


    it('should return all recent vasps with pending review status', () => {
        const vasps = [
            {
                "id": faker.datatype.uuid(),
                "name": faker.company.companyName,
                "common_name": faker.company.catchPhrase,
                "registered_directory": "trisatest.net",
                "verification_status": "VERIFIED",
                "last_updated": "2022-09-09T15:25:49Z",
                "verified_on": "2021-06-11T21:37:52Z",
                "traveler": false,
                "certificate_serial_number": "58482A7F46A7567EDC4B568E5829EAEE",
                "certificate_issued": "2020-11-23T21:50:17Z",
                "certificate_expiration": "2023-07-23T21:50:16Z",
                "verified_contacts": {
                    "administrative": true,
                    "billing": false,
                    "legal": true,
                    "technical": true
                }
            },
            {
                "id": faker.datatype.uuid(),
                "name": faker.company.companyName,
                "common_name": faker.company.catchPhrase,
                "registered_directory": "trisatest.net",
                "verification_status": "PENDING_REVIEW",
                "last_updated": "2022-09-09T15:25:49Z",
                "verified_on": "2022-05-04T22:00:31Z",
                "traveler": false,
                "certificate_serial_number": "2E7A293158A5ECC06169BB9FB2EEE9A1",
                "certificate_issued": "2021-11-23T13:35:47Z",
                "certificate_expiration": "2023-06-23T13:35:46Z",
                "verified_contacts": {
                    "legal": true,
                    "technical": true
                }
            },
        ]

        expect(getRecentVasps(vasps)).toHaveLength(1)
    })

    it('should return all recent vasps with certificate_issued date less than 30days', () => {
        const vasps = [
            {
                "id": faker.datatype.uuid(),
                "name": faker.company.companyName,
                "common_name": faker.company.catchPhrase,
                "registered_directory": "trisatest.net",
                "verification_status": "VERIFIED",
                "last_updated": "2021-09-09T15:25:49Z",
                "verified_on": "2021-06-11T21:37:52Z",
                "traveler": false,
                "certificate_serial_number": "58482A7F46A7567EDC4B568E5829EAEE",
                "certificate_issued": "2021-10-23T21:50:17Z",
                "certificate_expiration": "2023-07-23T21:50:16Z",
                "verified_contacts": {
                    "administrative": true,
                    "billing": false,
                    "legal": true,
                    "technical": true
                }
            },
            {
                "id": faker.datatype.uuid(),
                "name": faker.company.companyName,
                "common_name": faker.company.catchPhrase,
                "registered_directory": "trisatest.net",
                "verification_status": "VERIFIED",
                "last_updated": "2022-09-09T15:25:49Z",
                "verified_on": "2022-05-04T22:00:31Z",
                "traveler": false,
                "certificate_serial_number": "2E7A293158A5ECC06169BB9FB2EEE9A1",
                "certificate_issued": "2021-11-20T13:35:47Z",
                "certificate_expiration": "2023-06-23T13:35:46Z",
                "verified_contacts": {
                    "legal": true,
                    "technical": true
                }
            },
        ]

        expect(getRecentVasps(vasps)).toHaveLength(1)
    })

    it('should return all recent vasps with certificate_expiration date less than 30days', () => {
        const vasps = [
            {
                "id": faker.datatype.uuid(),
                "name": faker.company.companyName,
                "common_name": faker.company.catchPhrase,
                "registered_directory": "trisatest.net",
                "verification_status": "VERIFIED",
                "last_updated": "2022-09-09T15:25:49Z",
                "verified_on": "2021-06-11T21:37:52Z",
                "traveler": false,
                "certificate_serial_number": "58482A7F46A7567EDC4B568E5829EAEE",
                "certificate_issued": "2020-10-23T21:50:17Z",
                "certificate_expiration": "2021-12-02T21:50:16Z",
                "verified_contacts": {
                    "administrative": true,
                    "billing": false,
                    "legal": true,
                    "technical": true
                }
            },
            {
                "id": faker.datatype.uuid(),
                "name": faker.company.companyName,
                "common_name": faker.company.catchPhrase,
                "registered_directory": "trisatest.net",
                "verification_status": "VERIFIED",
                "last_updated": "2022-09-09T15:25:49Z",
                "verified_on": "2022-05-04T22:00:31Z",
                "traveler": false,
                "certificate_serial_number": "2E7A293158A5ECC06169BB9FB2EEE9A1",
                "certificate_issued": "2020-11-23T13:35:47Z",
                "certificate_expiration": "2023-06-23T13:35:46Z",
                "verified_contacts": {
                    "legal": true,
                    "technical": true
                }
            },
        ]

        expect(getRecentVasps(vasps)).toHaveLength(1)
    })
})