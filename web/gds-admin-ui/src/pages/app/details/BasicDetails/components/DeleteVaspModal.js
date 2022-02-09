import DeleteModal from 'components/DeleteModal';
import React from 'react';
import toast from 'react-hot-toast';
import { useHistory } from 'react-router-dom';
import { useParams } from 'react-router-dom';
import { deleteVasp } from 'services/vasp';

function DeleteVaspModal() {
    const [isLoading, setIsLoading] = React.useState(false)
    const params = useParams()
    const history = useHistory()

    const handleDeleteClick = () => {
        setIsLoading(true)
        if (params && params.id) {

            deleteVasp(params.id).then(() => {
                setIsLoading(false)
                history.replace('/vasps')


                toast.success("Registration deleted successfully", { duration: 7000 })
            }).catch(error => {
                toast.error("Unable to delete this registration")
                console.error('[deleteVasp] error', error)
                setIsLoading(false)
            })
        }
    }

    return <DeleteModal onDelete={handleDeleteClick} isLoading={isLoading} />
}

export default DeleteVaspModal;
