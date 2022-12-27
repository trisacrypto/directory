import faker from 'faker'
import Contact from "../Contact"
import { VERIFIED_CONTACT_STATUS_LABEL } from "../../../../../constants"
import { Modal } from "components/Modal"
import { render, screen } from 'utils/test-utils'


describe('Contact', () => {
    let verifiedContact
    let type
    let data
    let status;
    let contactNode

    it('should be verified', () => {
        verifiedContact = { legal: "Vita_VonRueden89@hotmail.com", technical: "Vita_VonRueden89@hotmail.com" }
        type = faker.random.arrayElement(["legal", "technical"])
        data = {
            email: 'Cielo.Kemmer67@yahoo.com',
            extra: null,
            name: 'Garland Goodwin',
            person: null,
            phone: '(551) 777-6790 x9018'
        };

        render(
            <Modal>
                <Contact verifiedContact={verifiedContact} type={type} data={data} />
            </Modal>
        )
        status = screen.getByTestId('verifiedContactStatus')
        contactNode = screen.getByTestId('contact-node')

        expect(status.textContent).toBe(VERIFIED_CONTACT_STATUS_LABEL.VERIFIED)
        expect(contactNode).not.toHaveClass()
    })

    it('should be alternate verified', () => {
        verifiedContact = { administrative: "Ozella_Crooks25@yahoo.com", billing: "Vita_VonRueden89@hotmail.com" }
        type = faker.random.arrayElement(["legal", "technical"])
        data = {
            email: 'Ozella_Crooks25@yahoo.com',
            extra: null,
            name: 'Julie Lowe',
            person: null,
            phone: '(827) 631-9433 x326'
        }

        render(
            <Modal>
                <Contact verifiedContact={verifiedContact} type={type} data={data} />
            </Modal>
        )

        status = screen.getByTestId('verifiedContactStatus')
        contactNode = screen.getByTestId('contact-node')

        expect(status.textContent).toBe(VERIFIED_CONTACT_STATUS_LABEL.ALTERNATE_VERIFIED)
    })

    it('should be unverified', () => {
        verifiedContact = { administrative: "VonRueden89@hotmail.com", billing: "Vita_VonRueden89@hotmail.com" }
        type = faker.random.arrayElement(["legal", "technical"])
        data = {
            email: 'Alia.Stehr45@gmail.com',
            extra: null,
            name: 'Kirk Bins',
            person: null,
            phone: '1-935-214-3799 x881'
        }

        render(
            <Contact verifiedContact={verifiedContact} type={type} data={data} />
        )

        status = screen.getByTestId('verifiedContactStatus')
        contactNode = screen.getByTestId('contact-node')

        expect(status.textContent).toBe('')
    })
})
