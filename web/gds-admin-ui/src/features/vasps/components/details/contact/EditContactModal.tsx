import { Contact } from '../../../types/Contact';
import EditContactForm from './EditContactForm';

const formDescriptionsText = {
    administrative:
        'Administrative or executive contact for your organization to field high-level requests or queries. (Strongly recommended).',
    technical:
        'Primary contact for handling technical queries about the operation and status of your service participating in the TRISA network. Can be a group or admin email. (Required).',
    legal: 'Compliance officer or legal contact for requests about the compliance requirements and legal status of your organization. (Required).',
    billing:
        'Billing contact for your organization to handle account and invoice requests or queries relating to the operation of the TRISA network. (Optional).',
};

type EditContactModalProps = {
    contactType: keyof typeof formDescriptionsText;
    contact: Contact;
};

function EditContactModal({ contactType, contact }: EditContactModalProps) {
    return (
        <>
            <h3 className="header-title mb-1">
                <span className="text-capitalize">{contactType}</span> Contacts
            </h3>
            <p className="text-muted font-16 mb-3">
                Please supply contact information for representatives of your organization. All contacts will receive an
                email verification token and the contact email must be verified before the registration can proceed.
            </p>
            <p className="text-muted font-14 mb-1">{formDescriptionsText[contactType]}</p>
            <EditContactForm contactType={contactType} contact={contact} />
        </>
    );
}

export default EditContactModal;
