import React from 'react';
import {
  Tr,
  Box,
  Text,
  Thead,
  Table,
  Flex,
  Tbody,
  Th,
  Stack,
  Heading,
  useColorModeValue,
  useColorMode
} from '@chakra-ui/react';
import CertificateRegistrationRow from 'components/Tables/CertificateRegistrationRow';
import Card from 'components/ui/Card';
import { t } from '@lingui/macro';
import { Trans } from '@lingui/react';
const defaultRowData = [
  {
    section: '1',
    name: t`Business Details`,
    description: t`Website, incorporation Date, VASP Category`,
    status: null
  },
  {
    section: '2',
    name: t`Legal Person`,
    description: t`Name, Addresss, Country, National Identifier`,
    status: null
  },
  {
    section: '3',
    name: t`Contacts`,
    description: t`Compliance, Technical, Admininstrative, Billing`,
    status: null
  },
  {
    section: '4',
    name: t`TRISA Implementation`,
    description: t`TRISA endpoint for communication`,
    status: null
  },
  {
    section: '5',
    name: t`TRIXO Questionnaire`,
    description: t`CDD and data protection policies`,
    status: null
  },
  {
    section: '6',
    name: t`Submit`,
    description: t`Final review and form submission`,
    status: null
  }
];
const CertificateRegistration = ({ data }: any) => {
  const textColor = useColorModeValue('#858585', 'white');

  return (
    <Card>
      <Stack p={4} mb={5}>
        <Heading fontSize="20px" fontWeight="bold" pb=".5rem">
          <Trans id="Certificate Registration Process">Certificate Registration Process</Trans>
        </Heading>
      </Stack>
      <Stack px={'60px'}>
        <Table
          color={textColor}
          width={'100%'}
          sx={{
            borderCollapse: 'separate',
            borderSpacing: '0 10px',
            Th: {
              textTransform: 'capitalize',
              color: '#858585',
              fontWeight: 'bold',
              borderBottom: 'none',
              fontSize: '0.9rem',
              textAlign: 'center'
            }
          }}>
          <Thead>
            <Tr>
              <Th>
                <Trans id="Section">Section</Trans>
              </Th>
              <Th>
                <Trans id="Name">Name</Trans>
              </Th>
              <Th>
                <Trans id="Description">Description</Trans>
              </Th>
              <Th>
                <Trans id="Status">Status</Trans>
              </Th>
              <Th>
                <Trans id="Action">Action</Trans>
              </Th>
            </Tr>
          </Thead>
          <Tbody>
            {defaultRowData.map((row: any) => {
              return (
                <CertificateRegistrationRow
                  key={row.section}
                  section={row.section}
                  name={row.name}
                  description={row.description}
                  status={row.status}
                />
              );
            })}
          </Tbody>
        </Table>
      </Stack>
    </Card>
  );
};

export default CertificateRegistration;
