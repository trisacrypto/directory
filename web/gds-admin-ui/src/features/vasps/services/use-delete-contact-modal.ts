import { useModalContext } from '@/components/Modal';
import { useParams } from 'react-router-dom';
import { useDeleteContact } from './delete-contact';
import { captureException } from '@sentry/react';
import toast from 'react-hot-toast';
import { queryClient } from '@/lib/react-query';

export default function useDeleteContactModal({ contactType }: { contactType: string }) {
    const params = useParams<{ id: string }>();
    const { closeModal } = useModalContext();
    const { mutate: deleteContact } = useDeleteContact();

    const handleDeleteClick = () => {
        deleteContact(
            {
                vaspId: params.id,
                kind: contactType,
            },
            {
                onSuccess() {
                    closeModal();
                    queryClient.invalidateQueries(['get-vasps', params.id]);
                },
                onError(error) {
                    toast.error('Could not delete this contact. Please try later!');
                    captureException(error);
                },
            }
        );
    };

    return {
        handleDeleteClick,
    };
}
