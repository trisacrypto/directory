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
            {
                "id": "00634f29-e22e-48f5-be2a-74feeee33464",
                "name": "Net Marketing Services Pte Ltd",
                "common_name": "eczesgsg.trisa.test-travelrule.sygna.io",
                "registered_directory": "trisatest.net",
                "verification_status": "VERIFIED",
                "last_updated": "2022-09-18T18:51:40Z",
                "verified_on": "2022-02-10T17:12:23Z",
                "traveler": false,
                "certificate_serial_number": "4BB35A55DF33B7561411DC1D7DCEFE32",
                "certificate_issued": "2022-03-18T21:14:40Z",
                "certificate_expiration": "2023-04-18T21:14:39Z",
                "verified_contacts": {
                    "administrative": false,
                    "legal": true,
                    "technical": false
                }
            },
            {
                "id": "033cd7d8-747e-455e-8565-be78538df1bf",
                "name": "Apex Crypto",
                "common_name": "trisa-9ffb8b3247f15fb585c98ea982f1d5a7.traveler.ciphertrace.com",
                "registered_directory": "trisatest.net",
                "verification_status": "VERIFIED",
                "last_updated": "2022-09-18T18:51:40Z",
                "verified_on": "2021-12-07T20:22:00Z",
                "traveler": true,
                "certificate_serial_number": "7385C056E92D60F38E347DE5735614BA",
                "certificate_issued": "2021-12-07T20:31:15Z",
                "certificate_expiration": "2023-01-07T20:31:14Z",
                "verified_contacts": {
                    "administrative": true,
                    "billing": true,
                    "legal": true,
                    "technical": true
                }
            },
            {
                "id": "03faf7d2-451d-4d90-8302-e80f0cc9848a",
                "name": "Guidehouse Inc.",
                "common_name": "trisa-a8c416cc67ea62e8cc30a2abebe066db.traveler.ciphertrace.com",
                "registered_directory": "trisatest.net",
                "verification_status": "PENDING_REVIEW",
                "last_updated": "2022-09-18T18:51:40Z",
                "verified_on": "2021-07-29T19:11:03Z",
                "traveler": true,
                "certificate_serial_number": "4AD8582453937C255ED1A96B5796B65A",
                "certificate_issued": "2021-08-05T20:57:45Z",
                "certificate_expiration": "2022-09-05T20:57:44Z",
                "verified_contacts": {
                    "legal": true,
                    "technical": true
                }
            },
            {
                "id": "063daf09-e5cd-4daa-b4dd-e2a800bdd678",
                "name": "CipherTrace Inc",
                "common_name": "trisa-e130acb71908fc8a87c5b1cd38ff2ade.traveler.ciphertrace.com",
                "registered_directory": "trisatest.net",
                "verification_status": "PENDING_REVIEW",
                "last_updated": "2022-09-18T18:51:40Z",
                "verified_on": "2021-06-23T17:46:10Z",
                "traveler": true,
                "certificate_serial_number": "2A50A7CD20A726D0F05484B518B65727",
                "certificate_issued": "2022-07-13T00:30:25Z",
                "certificate_expiration": "2023-08-13T00:30:24Z",
                "verified_contacts": {
                    "administrative": true,
                    "billing": false,
                    "legal": false,
                    "technical": true
                }
            },
            {
                "id": "07c102b3-5171-4d42-b976-ec001d4b8095",
                "name": "CoinFlip",
                "common_name": "trisa-1b640ffb0cb36e1cea44445969d18898.traveler.ciphertrace.com",
                "registered_directory": "trisatest.net",
                "verification_status": "PENDING_REVIEW",
                "last_updated": "2021-09-18T18:51:40Z",
                "verified_on": "2021-09-14T10:40:30Z",
                "traveler": true,
                "certificate_serial_number": "53632ED00AE5DA66663B4F0A85FB887A",
                "certificate_issued": "2021-09-14T10:48:50Z",
                "certificate_expiration": "2022-10-14T10:48:49Z",
                "verified_contacts": {
                    "administrative": true,
                    "billing": true,
                    "legal": true,
                    "technical": true
                }
            }
        ]

        expect(getRecentVasps(vasps)).toHaveLength(5)
    })
})