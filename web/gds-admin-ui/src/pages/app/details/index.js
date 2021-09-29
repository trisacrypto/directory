// @flow
import React from 'react';
import { Row, Col, } from 'react-bootstrap';

import PageTitle from '../../../components/PageTitle';
import { useParams } from "react-router-dom"

import Contact from './contact';
import BasicDetails from './BasicDetails';
import CertificateDetails from './CertificateDetails';
import { getVasp } from "../../../services/vasps"
import TrixoForm from './TrixoForm';
import Ivms from './ivms';

const VaspDetails = (): React$Element<React$FragmentType> => {
    const [vasp, setVasp] = React.useState({});
    const params = useParams();

    React.useEffect(() => {
        if (params && params.id) {
            getVasp(params.id).then(response => {
                setVasp(response.data)
            }).catch(error => {
                console.error("[BasicDetails] getVasp", error.message)
            })
        }
    }, [params])


    return (
        <React.Fragment>
            <PageTitle
                breadCrumbItems={[
                    { label: 'List', path: '/vasps' },
                    { label: 'Details', path: `/vasps/${params?.id}`, active: true }
                ]}
                title={'Registration Details'}
            />
            {vasp && (
                <>
                    <Row>
                        <Col>
                            <BasicDetails data={vasp} />
                        </Col>
                        <Col>
                            <Contact data={vasp?.vasp?.contacts} />
                        </Col>
                    </Row>
                    <Row>
                        <Col md={6}>
                            <Ivms data={vasp.vasp?.entity} />
                            <TrixoForm data={vasp.trixo} />
                        </Col>
                        <Col md={6}>
                            <CertificateDetails data={vasp?.vasp?.identity_certificate} />
                        </Col>
                    </Row>
                </>
            )}
        </React.Fragment>
    );
};

export default VaspDetails;
