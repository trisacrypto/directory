import { Dropdown, Row } from 'react-bootstrap';

import { Modal, ModalContent, ModalOpenButton } from '@/components/Modal';
import { isOptionAvailable } from '@/utils';

import BusinessInfosForm from './BusinessInfosForm';
import DeleteVaspModal from './DeleteVaspModal';
import Ivms101RecordForm from './Ivms101RecordForm';
import Print from './Print';
import ReviewForm from './ReviewForm';
import TrisaImplementationDetailsForm from './TrisaImplementationDetailsForm';
import useGetBasicDetailsDropdown from '@/features/vasps/services/use-basic-details-dropdown';

const BasicDetailsDropDown = ({ isNotPendingReview, vasp }: any) => {
    const { closeEmailDrawer, handlePrint } = useGetBasicDetailsDropdown({ vasp });

    return (
        <Dropdown className="float-end" align="end">
            <Dropdown.Toggle
                data-testid="dripicons-dots-3"
                variant="link"
                as="a"
                className="card-drop arrow-none cursor-pointer p-0 shadow-none">
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
                    <Modal>
                        <ModalOpenButton>
                            <Dropdown.Item disabled={!isOptionAvailable(vasp?.vasp?.verification_status)}>
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

                <Print onPrint={handlePrint} />
                <Dropdown.Item onClick={closeEmailDrawer}>
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
