import { useModal } from '@/contexts/modal';
import { useGetReviewNotes } from './get-review-notes';
import slugify from 'slugify';
import { VaspDocument } from '../components';
import { downloadFile } from '@/helpers/api/utils';
import nProgress from 'nprogress';
import { captureException } from '@sentry/react';
import toast from 'react-hot-toast';
import { pdf } from '@react-pdf/renderer';

export default function useGetBasicDetailsDropdown({ vasp }: any) {
    const { openSendEmailModal } = useModal();
    const { data: reviewNotes } = useGetReviewNotes({ vaspId: vasp?.vasp?.id });

    const handleClose = () => openSendEmailModal({ name: vasp?.name, id: vasp?.vasp?.id });

    const getFilename = () => `${Date.now()}-${slugify(vasp?.name || '')}`;

    const generatePdfDocument = async () => {
        nProgress.start();
        try {
            const blob = await pdf(<VaspDocument vasp={vasp} notes={reviewNotes} />).toBlob();
            downloadFile(blob, `${getFilename()}.pdf`, 'application/pdf');
            nProgress.done();
        } catch (error) {
            captureException(error);
            toast.error('Unable to export as PDF');
            nProgress.done();
        }
    };

    return {
        handlePrint: generatePdfDocument,
        closeEmailDrawer: handleClose,
    };
}
