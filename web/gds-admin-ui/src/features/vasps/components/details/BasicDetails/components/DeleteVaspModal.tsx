import React from 'react';
import toast from 'react-hot-toast';
import { useHistory, useParams } from 'react-router-dom';

import DeleteModal from '@/components/DeleteModal';
import { deleteVasp } from '@/services/vasp';
import { useGetVasp } from '@/features/vasps/services';
import { captureException } from '@sentry/react';

function DeleteVaspModal() {
    const [isLoading, setIsLoading] = React.useState(false);
    const params = useParams<{ id: string }>();
    const { data: vasp } = useGetVasp({ vaspId: params.id });
    const history = useHistory();

    const handleDeleteClick = () => {
        setIsLoading(true);
        if (params && params.id) {
            deleteVasp(params.id)
                .then(() => {
                    setIsLoading(false);
                    history.replace('/vasps');

                    toast.success('Registration deleted successfully', { duration: 7000 });
                })
                .catch((error) => {
                    toast.error('Unable to delete this registration');
                    captureException(error);
                    setIsLoading(false);
                });
        }
    };

    return <DeleteModal onDelete={handleDeleteClick} isLoading={isLoading} vaspId={params.id} vasp={vasp} />;
}

export default DeleteVaspModal;
