// @flow
import React from 'react';
import { Row, Col, } from 'react-bootstrap';

import PageTitle from '../../../components/PageTitle';
import { useParams } from "react-router-dom"

import ContactInfos from './ContactInfos';
import BasicDetails from './BasicDetails';
import CertificateDetails from './CertificateDetails';
import { getVasp } from "../../../services/vasps"
import TrixoForm from './TrixoForm';
import Ivms from './Ivms';

const VaspDetails = (): React$Element<React$FragmentType> => {
    const [vasp, setVasp] = React.useState({ });
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
                breadCrumbItems={[]}
                title={'VASP Details'}
            />
            <Row>
                <Col>
                    <BasicDetails data={vasp} />
                </Col>
            </Row>
            <Row>
                <Col>
                    <TrixoForm data={vasp.trixo} />
                </Col>
            </Row>
            <Row>
                <Col>
                    <Ivms data={vasp.entity} />
                </Col>
            </Row>
            <Row>
                <Col sm={12}>
                    <ContactInfos data={vasp.contacts} />
                </Col>
            </Row>
            <Row>
                <Col>
                    <CertificateDetails data={vasp.identity_certificate} />
                </Col>
            </Row>
        </React.Fragment>
    );
};

export default VaspDetails;
