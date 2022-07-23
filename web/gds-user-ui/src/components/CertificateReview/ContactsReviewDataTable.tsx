import { Stack, Table, Tbody, Tr, Td, Text } from '@chakra-ui/react';
import { t } from '@lingui/macro';
import React, { useEffect } from 'react';
import { useSelector, RootStateOrAny } from 'react-redux';
import { TStep, loadDefaultValueFromLocalStorage } from 'utils/localStorageHelper';
import { hasValue } from 'utils/utils';

function ContactsReviewDataTable() {
  const steps: TStep[] = useSelector((state: RootStateOrAny) => state.stepper.steps);
  const [contacts, setContacts] = React.useState<any>({});
  useEffect(() => {
    const getStepperData = loadDefaultValueFromLocalStorage();
    const stepData = {
      ...getStepperData.contacts
    };
    setContacts(stepData);
  }, [steps]);

  return (
    <Stack fontSize={'1rem'}>
      <Table
        sx={{
          'td:nth-child(2),td:nth-child(3)': { fontWeight: 'semibold' },
          Tr: { borderStyle: 'hidden' },
          'td:first-child': {
            width: '50%',
            paddingLeft: '1rem'
          },
          td: {
            borderBottom: 'none',
            paddingInlineStart: 0,
            paddingY: 2.5
          }
        }}>
        <Tbody>
          {['technical', 'legal', 'administrative', 'billing'].map((contact, index) => (
            <Tr key={index}>
              <Td textTransform="capitalize">
                {t`${contact === 'legal' ? `Compliance / ${contact}` : contact} Contact`}
              </Td>
              <Td>
                {hasValue(contacts?.[contact]) ? (
                  <>
                    {contacts?.[contact]?.name && <Text>{contacts?.[contact]?.name}</Text>}
                    {contacts?.[contact]?.email && <Text>{contacts?.[contact]?.email}</Text>}
                    {contacts?.[contact]?.phone && <Text>{contacts?.[contact]?.phone}</Text>}
                  </>
                ) : (
                  'N/A'
                )}
              </Td>
            </Tr>
          ))}
        </Tbody>
      </Table>
    </Stack>
  );
}

export default ContactsReviewDataTable;
