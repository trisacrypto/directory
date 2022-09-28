import { Stack, Table, Tbody, Tr, Td, Text } from '@chakra-ui/react';
import { t } from '@lingui/macro';
import { hasValue } from 'utils/utils';
interface ContactsProps {
  data?: any;
}
function ContactsReviewDataTable({ data }: ContactsProps) {
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
          {['legal', 'technical', 'administrative', 'billing'].map((contact, index) => (
            <Tr key={index}>
              <Td textTransform="capitalize">
                {t`${contact === 'legal' ? `Compliance / ${contact}` : contact} Contact`}
              </Td>
              <Td>
                {hasValue(data?.[contact]) ? (
                  <>
                    {data?.[contact]?.name && <Text>{data?.[contact]?.name}</Text>}
                    {data?.[contact]?.email && <Text>{data?.[contact]?.email}</Text>}
                    {data?.[contact]?.phone && <Text>{data?.[contact]?.phone}</Text>}
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
