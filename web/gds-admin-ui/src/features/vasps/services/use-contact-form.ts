import { useModalContext } from '@/components/Modal';
import { getContactInitialValues } from '@/utils/form-references';
import { useForm } from 'react-hook-form';
import { useUpdateContact } from './update-contact';
import { useParams } from 'react-router-dom';

export default function useContactForm({ contactType, contact }: any) {
    const hookformMethods = useForm({
        defaultValues: getContactInitialValues(contact),
        mode: 'onChange',
    });
    const params = useParams<{ id: string }>();
    const { closeModal } = useModalContext();
    const { mutate: updateContact, isError, error } = useUpdateContact();

    const onSubmit = (data: any) => {
        if (params && params.id) {
            updateContact(
                {
                    vaspId: params.id,
                    kind: contactType,
                    data,
                },
                {
                    onSuccess() {
                        closeModal();
                    },
                }
            );
        }
    };

    const handleAlertClose = () => {};

    return {
        handleAlertClose,
        onSubmit,
        isError,
        error,
        ...hookformMethods,
    };
}
