import PropTypes from 'prop-types';
import { Dropdown, Row } from 'react-bootstrap';

import { Modal, ModalContent, ModalOpenButton } from '@/components/Modal';

import TrixoEditForm from './TrixoEditForm';

function TrixoDropdown({ data }: { data: any }) {
    return (
        <div>
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
                            <Dropdown.Item>
                                <i className="mdi mdi-square-edit-outline me-1" />
                                Edit
                            </Dropdown.Item>
                        </ModalOpenButton>
                        <ModalContent>
                            <Row className="p-4">
                                <h3>Edit TRIXO Form</h3>
                                <p>
                                    This questionnaire is designed to help the TRISA working group and TRISA members
                                    understand the regulatory regime of your organization. The information you provide
                                    will help ensure that required compliance information exchanges are conducted
                                    correctly and safely. All verified TRISA members will have access to this
                                    information.
                                </p>
                                <TrixoEditForm data={data} />
                            </Row>
                        </ModalContent>
                    </Modal>
                </Dropdown.Menu>
            </Dropdown>
        </div>
    );
}

TrixoDropdown.propTypes = {
    handleEditClick: PropTypes.func,
};

export default TrixoDropdown;
