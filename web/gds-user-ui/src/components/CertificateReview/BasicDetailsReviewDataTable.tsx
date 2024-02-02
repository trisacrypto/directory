import { Stack, Table, Tbody, Tr, Td, Tag, Link } from '@chakra-ui/react';
import { Trans } from '@lingui/react';
import { BUSINESS_CATEGORY, getBusinessCategiryLabel } from 'constants/basic-details';
interface BasicReviewProps {
  data?: any;
}
function BasicDetailsReviewDataTable({ data }: BasicReviewProps) {
  return (
    <Stack fontSize={'1rem'}>
      <Table
        sx={{
          'td:nth-child(2),td:nth-child(3)': { fontWeight: 'semibold' },
          'td:first-child': {
            width: '50%'
          },
          td: {
            borderBottom: 'none',
            paddingInlineStart: 0,
            paddingY: 2.5
          }
        }}>
        <Tbody>
          <Tr>
            <Td pl={'1rem !important'}>
              <Trans id="Organization Name">Organization Name</Trans>
            </Td>
            <Td>{data?.organization_name || 'N/A'}</Td>
            <Td></Td>
          </Tr>
          <Tr>
            <Td borderBottom={'none'} pl={'1rem !important'}>
              <Trans id="Website">Website</Trans>
            </Td>
            <Td borderBottom={'none'} whiteSpace="break-spaces" lineHeight={1.5}>
              {data.website ? (
                <Link href={data.website} isExternal>
                  {data.website}
                </Link>
              ) : (
                'N/A'
              )}
            </Td>
            <Td></Td>
          </Tr>
          <Tr>
            <Td pl={'1rem !important'}>
              <Trans id="Business Category">Business Category</Trans>
            </Td>
            <Td>{(BUSINESS_CATEGORY as any)[data.business_category] || 'N/A'}</Td>
            <Td></Td>
          </Tr>
          <Tr borderStyle={'hidden'}>
            <Td pl={'1rem !important'} whiteSpace="break-spaces" lineHeight={1.5}>
              <Trans id="Date of Incorporation / Establishment">
                Date of Incorporation / Establishment
              </Trans>
            </Td>
            <Td>{data.established_on || 'N/A'}</Td>
            <Td></Td>
          </Tr>
          <Tr borderStyle={'hidden'}>
            <Td pl={'1rem !important'} whiteSpace="break-spaces" lineHeight={1.5}>
              <Trans id="VASP Category">VASP Category</Trans>
            </Td>
            <Td>
              {data?.vasp_categories && data?.vasp_categories.length
                ? data?.vasp_categories?.map((categ: any) => {
                    return (
                      <Tag key={categ} color={'white'} bg={'blue'} mr={2} mb={1}>
                        {getBusinessCategiryLabel(categ)}
                      </Tag>
                    );
                  })
                : 'N/A'}
            </Td>
            <Td></Td>
          </Tr>
        </Tbody>
      </Table>
    </Stack>
  );
}

export default BasicDetailsReviewDataTable;
