import React from 'react';
import { Card, Dropdown, Table } from 'react-bootstrap';
import { useHistory } from 'react-router-dom'
import dayjs from 'dayjs';
import { StatusLabel } from 'constants/index';
import TrisaFavicon from 'assets/images/trisa_favicon.png'
import CiphertraceFavicon from 'assets/images/ciphertrace.ico'
import PropTypes from 'prop-types';

import relativeTime from 'dayjs/plugin/relativeTime'
import { actionType, useModal } from 'contexts/modal';
import OvalLoader from 'components/OvalLoader';
dayjs.extend(relativeTime)




const Tasks = ({ data, isLoading }) => {
    const [vasp, setVasp] = React.useState({ name: '', id: '' });
    const history = useHistory()
    const { dispatch } = useModal()

    const handleResendEmailClick = (name) => {
        setVasp(name)
        dispatch({ type: actionType.SEND_EMAIL_MODAL, payload: { vasp } })
    }

    return (
        <Card style={{ height: '95%' }}>
            <Card.Body>
                <h4 className="header-title mb-3">Pending Reviews</h4>
                {
                    isLoading ?
                        <OvalLoader />
                        :
                        <Table className="table table-centered table-nowrap table-hover mb-0 z-index-2">
                            <tbody>
                                {
                                    Array.isArray(data?.vasps) && data.vasps.map((vasp) => (
                                        <tr key={vasp.id}>
                                            <td onClick={() => history.push(`/vasps/${vasp?.id}`)} className="d-flex gap-2 align-items-center" role="button">
                                                <div>
                                                    {
                                                        vasp?.traveler ? <img src={CiphertraceFavicon} width="30" alt="Cyphertrace" /> : <img src={TrisaFavicon} width="30" className="img-fluid" alt="Trisa" />
                                                    }
                                                </div>
                                                <div>
                                                    <h5 className="font-14 my-1 d-flex gap-2">
                                                        {vasp?.name}
                                                    </h5>
                                                    <span className="text-muted font-13">{dayjs(vasp?.last_updated).fromNow()}</span>
                                                </div>
                                            </td>
                                            <td>
                                                <span className="text-muted font-13">Status</span> <br />
                                                <span className="badge badge-warning-lighten">{StatusLabel[vasp?.verification_status]}</span>
                                            </td>
                                            <td>
                                                <span className="text-muted font-13">Emails</span>
                                                <div className="font-14 mt-1 fw-normal">
                                                    {
                                                        vasp?.verified_contacts && Object.keys(vasp?.verified_contacts).map(contact => (
                                                            <span key={contact} className={`badge ${vasp?.verified_contacts[contact] ? "badge-success-lighten" : "badge-danger-lighten"}`} style={{ marginRight: 4 }}>{contact}</span>
                                                        )
                                                        )
                                                    }
                                                </div>
                                            </td>
                                            <td className="table-action text-center" style={{ width: '90px' }}>
                                                <Dropdown align="end">
                                                    <Dropdown.Toggle variant="link" className="arrow-none card-drop p-0 shadow-none">
                                                        <i className="mdi mdi-dots-horizontal"></i>
                                                    </Dropdown.Toggle>
                                                    <Dropdown.Menu>
                                                        <Dropdown.Item onClick={() => history.push(`/vasps/${vasp?.id}`)}> <span className="mdi mdi-eye-outline"></span> View</Dropdown.Item>
                                                        <Dropdown.Item onClick={() => handleResendEmailClick(vasp)}> <span className="mdi mdi-email-edit-outline"></span> Resend email</Dropdown.Item>
                                                    </Dropdown.Menu>
                                                </Dropdown>
                                            </td>
                                        </tr>
                                    ))
                                }
                            </tbody>
                        </Table>
                }
            </Card.Body>
        </Card>
    );
};

Tasks.propTypes = {
    data: PropTypes.object
}

export default Tasks;
