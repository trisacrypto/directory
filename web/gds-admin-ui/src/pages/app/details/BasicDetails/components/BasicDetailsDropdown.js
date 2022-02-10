
import { pdf } from '@react-pdf/renderer'
import { Modal, ModalContent, ModalOpenButton } from 'components/Modal'
import { actionType, useModal } from 'contexts/modal'
import { downloadFile } from 'helpers/api/utils'
import nProgress from 'nprogress'
import React from 'react'
import { Dropdown, Row } from 'react-bootstrap'
import { useSelector } from 'react-redux'
import { getAllReviewNotes } from 'redux/selectors'
import { isOptionAvailable } from 'utils'
import VaspDocument from '../../VaspDocument'
import BusinessInfosForm from './BusinessInfosForm'
import DeleteVaspModal from './DeleteVaspModal'
import Ivms101RecordForm from './Ivms101RecordForm'
import ReviewForm from './ReviewForm'
import TrisaImplementationDetailsForm from './TrisaImplementationDetailsForm'

const BasicDetailsDropDown = ({ isNotPendingReview, vasp }) => {
    const { dispatch } = useModal()
    const reviewNotes = useSelector(getAllReviewNotes)

    const handleClose = () => dispatch({ type: actionType.SEND_EMAIL_MODAL, payload: { vasp: { name: vasp?.name, id: vasp?.vasp?.id } } })

    const generatePdfDocument = async (filename) => {
        nProgress.start()
        try {
            const blob = await pdf(<VaspDocument vasp={vasp} notes={reviewNotes} />).toBlob()
            downloadFile(blob, `${filename}.pdf`, 'application/pdf')
            nProgress.done()
        } catch (error) {
            console.error('Unable to export as PDF', error)
            nProgress.done()
        }
    };

    return (
        <Dropdown className="float-end" align="end">
            <Dropdown.Toggle
                data-testid="dripicons-dots-3"
                variant="link"
                tag="a"
                className="card-drop arrow-none cursor-pointer p-0 shadow-none">
                <i className="dripicons-dots-3"></i>
            </Dropdown.Toggle>
            <Dropdown.Menu>
                <Modal>
                    <ModalOpenButton>
                        <Dropdown.Item data-testid="reviewItem" disabled={isNotPendingReview()}>
                            <i className="mdi mdi-card-search me-1"></i>Review
                        </Dropdown.Item>
                    </ModalOpenButton>
                    <ModalContent size="md">
                        <Row className='p-4'>
                            <ReviewForm />
                        </Row>
                    </ModalContent>
                </Modal >

                {
                    <>
                        <Modal>
                            <ModalOpenButton>
                                <Dropdown.Item>
                                    <i className="mdi mdi-briefcase-edit me-1"></i>Edit Business Info
                                </Dropdown.Item>
                            </ModalOpenButton>
                            <ModalContent size="lg">
                                <Row className='p-4'>
                                    <BusinessInfosForm data={vasp} />
                                </Row>
                            </ModalContent>
                        </Modal >
                        {
                           isOptionAvailable(vasp?.vasp?.verification_status) && (
                                <Modal>
                                    <ModalOpenButton>
                                        <Dropdown.Item>
                                            <i className="mdi mdi-network me-1"></i>Edit TRISA Details
                                        </Dropdown.Item>
                                    </ModalOpenButton>
                                    <ModalContent size="lg">
                                        <Row className='p-4'>
                                            <TrisaImplementationDetailsForm data={vasp} />
                                        </Row>
                                    </ModalContent>
                                </Modal >
                            )
                        }
                        <Modal>
                            <ModalOpenButton>
                                <Dropdown.Item>
                                    <i className="mdi mdi-office-building me-1"></i>Edit IVMS 101 Record
                                </Dropdown.Item>
                            </ModalOpenButton>
                            <ModalContent size="lg">
                                <Row className='p-4'>
                                    <Ivms101RecordForm data={vasp.vasp.entity} />
                                </Row>
                            </ModalContent>
                        </Modal >
                    </>
                }

                <Dropdown.Item onClick={() => generatePdfDocument(vasp?.name)}>
                    <i className="mdi mdi-printer me-1"></i>Print
                </Dropdown.Item>
                <Dropdown.Item onClick={handleClose}>
                    <i className="mdi mdi-email me-1"></i>Resend
                </Dropdown.Item>
                <Modal>
                    <ModalOpenButton>
                        <Dropdown.Item disabled={!isOptionAvailable(vasp?.vasp?.verification_status)}>
                            <i className="mdi mdi-trash-can me-1"></i>Delete
                        </Dropdown.Item>
                    </ModalOpenButton>
                    <ModalContent size="md">
                        <Row className='p-4'>
                            <DeleteVaspModal />
                        </Row>
                    </ModalContent>
                </Modal >
            </Dropdown.Menu >
        </Dropdown >
    )
}


export default BasicDetailsDropDown
