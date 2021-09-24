import React from 'react';
import { Card, Dropdown, Table } from 'react-bootstrap';
import { useSelector } from 'react-redux';
import { useHistory } from 'react-router-dom'
import dayjs from 'dayjs';
import ResendEmail from '../../../components/ResendEmail';
import { StatusLabel } from '../../../constants';
import TrisaFavicon from '../../../assets/images/trisa_favicon.png'
import CiphertraceFavicon from '../../../assets/images/ciphertrace.ico'

import relativeTime from 'dayjs/plugin/relativeTime'
dayjs.extend(relativeTime)




const Tasks = () => {
    const { vasps } = useSelector((state) => ({
        vasps: state.Vasps.data,
        certificates: state.Certificates.data
    }));
    const [modal, setModal] = React.useState(false);
    const [vaspName, setVaspName] = React.useState("");
    const history = useHistory()

    const toggle = () => {
        setModal(!modal);
    };

    const handleResendEmailClick = (name) => {
        setVaspName(name)
        toggle()
    }


    return (
        <Card>
            <Card.Body>
                <h4 className="header-title mb-3">Pending Reviews</h4>

                <Table responsive className="table table-centered table-nowrap table-hover mb-0 z-index-2">
                    <tbody>
                        {
                            Array.isArray(vasps?.vasps) && vasps?.vasps.map((vasp) => vasp?.verification_status === "PENDING_REVIEW" && (

                                <tr key={vasp.id}>
                                    <td className="d-flex gap-2 align-items-center">
                                        <div>
                                            {
                                                vasp?.is_traveler ? <img src={CiphertraceFavicon} className="img-fluid" alt="Cyphertrace" /> : <img src={TrisaFavicon} className="img-fluid" alt="Trisa" />
                                            }
                                        </div>
                                        <div>
                                            <h5 className="font-14 my-1 d-flex gap-2">
                                                <a href="/" className="text-body">
                                                    {vasp?.name}
                                                </a>
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
                                                <Dropdown.Item onClick={() => history.push(`/vasps/${vasp?.id}/details`)}> <span className="mdi mdi-eye-outline"></span> View</Dropdown.Item>
                                                <Dropdown.Item onClick={() => handleResendEmailClick(vasp?.name)}> <span className="mdi mdi-email-edit-outline"></span> Resend email</Dropdown.Item>
                                            </Dropdown.Menu>
                                        </Dropdown>
                                    </td>
                                </tr>
                            ))
                        }
                    </tbody>
                </Table>
                <ResendEmail toggle={toggle} modal={modal} vaspName={vaspName} />
            </Card.Body>
        </Card>
    );
};

export default Tasks;
