import { Stack, Table, Thead, Tr, Th, Tbody, useColorModeValue } from '@chakra-ui/react';
import { t } from '@lingui/macro';
import { Trans } from '@lingui/react';
import CertificateRegistrationRow from 'components/Tables/CertificateRegistrationRow';

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

function CertificateRegistrationTable() {
  const textColor = useColorModeValue('#858585', 'white');

  return (
    <Stack px={'60px'}>
      <Table color={textColor} width={'100%'} variant="simple">
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
  );
}

export default CertificateRegistrationTable;
