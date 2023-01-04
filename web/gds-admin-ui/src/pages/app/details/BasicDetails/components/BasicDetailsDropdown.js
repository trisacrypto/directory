import { pdf } from '@react-pdf/renderer';
import nProgress from 'nprogress';
import { Dropdown, Row } from 'react-bootstrap';
import toast from 'react-hot-toast';
import { useSelector } from 'react-redux';
import slugify from 'slugify';

import { Modal, ModalContent, ModalOpenButton } from '@/components/Modal';
import { actionType, useModal } from '@/contexts/modal';
import { downloadFile } from '@/helpers/api/utils';
import { getAllReviewNotes } from '@/redux/selectors';
import { isOptionAvailable } from '@/utils';

import VaspDocument from '../../VaspDocument';

import BusinessInfosForm from './BusinessInfosForm';
import DeleteVaspModal from './DeleteVaspModal';
import Ivms101RecordForm from './Ivms101RecordForm';
import Print from './Print';
import ReviewForm from './ReviewForm';
import TrisaImplementationDetailsForm from './TrisaImplementationDetailsForm';

const BasicDetailsDropDown = ({ isNotPendingReview, vasp }) => {
  const { dispatch } = useModal();
  const reviewNotes = useSelector(getAllReviewNotes);

  const handleClose = () =>
    dispatch({
      type: actionType.SEND_EMAIL_MODAL,
      payload: { vasp: { name: vasp?.name, id: vasp?.vasp?.id } },
    });

  const generatePdfDocument = async () => {
    const filename = `${Date.now()}-${slugify(vasp?.name || '')}`;
    nProgress.start();
    try {
      const blob = await pdf(<VaspDocument vasp={vasp} notes={reviewNotes} />).toBlob();
      downloadFile(blob, `${filename}.pdf`, 'application/pdf');
      nProgress.done();
    } catch (error) {
      // TODO: catch this error with sentry
      toast.error('Unable to export as PDF');
      nProgress.done();
    }
  };

  return (
    <Dropdown className="float-end" align="end">
      <Dropdown.Toggle
        data-testid="dripicons-dots-3"
        variant="link"
        tag="a"
        className="card-drop arrow-none cursor-pointer p-0 shadow-none"
      >
        <i className="dripicons-dots-3" />
      </Dropdown.Toggle>
      <Dropdown.Menu>
        <Modal>
          <ModalOpenButton>
            <Dropdown.Item data-testid="reviewItem" disabled={isNotPendingReview()}>
              <i className="mdi mdi-card-search me-1" />
              Review
            </Dropdown.Item>
          </ModalOpenButton>
          <ModalContent size="md">
            <Row className="p-4">
              <ReviewForm />
            </Row>
          </ModalContent>
        </Modal>

        <>
          <Modal>
            <ModalOpenButton>
              <Dropdown.Item>
                <i className="mdi mdi-briefcase-edit me-1" />
                Edit Business Info
              </Dropdown.Item>
            </ModalOpenButton>
            <ModalContent size="lg">
              <Row className="p-4">
                <BusinessInfosForm data={vasp} />
              </Row>
            </ModalContent>
          </Modal>
          {isOptionAvailable(vasp?.vasp?.verification_status) && (
            <Modal>
              <ModalOpenButton>
                <Dropdown.Item>
                  <i className="mdi mdi-network me-1" />
                  Edit TRISA Details
                </Dropdown.Item>
              </ModalOpenButton>
              <ModalContent size="lg">
                <Row className="p-4">
                  <TrisaImplementationDetailsForm data={vasp} />
                </Row>
              </ModalContent>
            </Modal>
          )}
          <Modal>
            <ModalOpenButton>
              <Dropdown.Item>
                <i className="mdi mdi-office-building me-1" />
                Edit IVMS 101 Record
              </Dropdown.Item>
            </ModalOpenButton>
            <ModalContent size="lg">
              <Row className="p-4">
                <Ivms101RecordForm data={vasp.vasp?.entity} />
              </Row>
            </ModalContent>
          </Modal>
        </>

        <Print onPrint={generatePdfDocument} />
        <Dropdown.Item onClick={handleClose}>
          <i className="mdi mdi-email me-1" />
          Resend
        </Dropdown.Item>
        <Modal>
          <ModalOpenButton>
            <Dropdown.Item disabled={!isOptionAvailable(vasp?.vasp?.verification_status)}>
              <i className="mdi mdi-trash-can me-1" />
              Delete
            </Dropdown.Item>
          </ModalOpenButton>
          <ModalContent size="md">
            <Row className="p-4">
              <DeleteVaspModal />
            </Row>
          </ModalContent>
        </Modal>
      </Dropdown.Menu>
    </Dropdown>
  );
};

export default BasicDetailsDropDown;
