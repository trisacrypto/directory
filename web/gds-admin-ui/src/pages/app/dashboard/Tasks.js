// @flow
import React from 'react';
import { Card, Dropdown, Table } from 'react-bootstrap';
import { useSelector } from 'react-redux';
import { Status } from '../../../constants';



const Tasks = (): React$Element<any> => {
    const { vasps } = useSelector((state) => ({
        vasps: state.Vasps.data,
        certificates: state.Certificates.data
    }));


    return (
        <Card>
            <Card.Body>
                <h4 className="header-title mb-3">Summary - Pending</h4>

                <Table responsive className="table table-centered table-nowrap table-hover mb-0">
                    <tbody>
                        {
                            Array.isArray(vasps) && vasps.map((vasp) => vasp?.verification_status === "PENDING_REVIEW" && (

                                <tr key={vasp.id}>
                                    <td>
                                        <h5 className="font-14 my-1">
                                            <a href="/" className="text-body">
                                                {vasp?.common_name}
                                            </a>
                                        </h5>
                                        <span className="text-muted font-13">{vasp?.established_on}</span>
                                    </td>
                                    <td>
                                        <span className="text-muted font-13">Status</span> <br />
                                        <span className="badge badge-warning-lighten">{Status[vasp?.verification_status]}</span>
                                    </td>
                                    <td>
                                        <span className="text-muted font-13">Percent verified</span>
                                        <h5 className="font-14 mt-1 fw-normal">1/4</h5>
                                    </td>
                                    <td className="table-action text-center" style={{ width: '90px' }}>
                                        <Dropdown align="end">
                                            <Dropdown.Toggle variant="link" className="arrow-none card-drop p-0 shadow-none">
                                                <i className="mdi mdi-dots-horizontal"></i>
                                            </Dropdown.Toggle>
                                            <Dropdown.Menu>
                                                <Dropdown.Item>Resend email</Dropdown.Item>
                                                <Dropdown.Item>Pending</Dropdown.Item>
                                                <Dropdown.Item>Review</Dropdown.Item>
                                            </Dropdown.Menu>
                                        </Dropdown>
                                    </td>
                                </tr>
                            ))
                        }
                    </tbody>
                </Table>
            </Card.Body>
        </Card>
    );
};

export default Tasks;
