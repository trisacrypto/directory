import { Stack, Table, Tbody, Tr, Td, Tag, Link } from '@chakra-ui/react';
import { Trans } from '@lingui/react';
import { BUSINESS_CATEGORY, getBusinessCategiryLabel } from 'constants/basic-details';
import React, { useEffect } from 'react';
import { RootStateOrAny, useSelector } from 'react-redux';
import { loadDefaultValueFromLocalStorage, TStep } from 'utils/localStorageHelper';

function BasicDetailsReviewbasicDetailTable() {
  const steps: TStep[] = useSelector((state: RootStateOrAny) => state.stepper.steps);
  const [basicDetail, setBasicDetail] = React.useState<any>({});

  useEffect(() => {
    const getStepperbasicDetail = loadDefaultValueFromLocalStorage();
    const stepbasicDetail = {
      website: getStepperbasicDetail.website,
      established_on: getStepperbasicDetail.established_on,
      vasp_categories: getStepperbasicDetail.vasp_categories,
      business_category: getStepperbasicDetail.business_category
    };

    setBasicDetail(stepbasicDetail);
  }, [steps]);

  return (
    <Stack fontSize={18}>
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
        <Tbody
          sx={{
            '*': {
              fontSize: '1rem'
            }
          }}>
          <Tr>
            <Td borderBottom={'none'} pl={'1rem !important'}>
              <Trans id="Website">Website</Trans>
            </Td>
            <Td borderBottom={'none'} whiteSpace="break-spaces" lineHeight={1.5}>
              {basicDetail.website ? (
                <Link href={basicDetail.website} isExternal>
                  {basicDetail.website}
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
            <Td>{(BUSINESS_CATEGORY as any)[basicDetail.business_category] || 'N/A'}</Td>
            <Td></Td>
          </Tr>
          <Tr borderStyle={'hidden'}>
            <Td pl={'1rem !important'} whiteSpace="break-spaces" lineHeight={1.5}>
              <Trans id="Date of Incorporation / Establishment">
                Date of Incorporation / Establishment
              </Trans>
            </Td>
            <Td>{basicDetail.established_on || 'N/A'}</Td>
            <Td></Td>
          </Tr>
          <Tr borderStyle={'hidden'}>
            <Td pl={'1rem !important'} whiteSpace="break-spaces" lineHeight={1.5}>
              <Trans id="VASP Category">VASP Category</Trans>
            </Td>
            <Td>
              {basicDetail?.vasp_categories && basicDetail?.vasp_categories.length
                ? basicDetail?.vasp_categories?.map((categ: any) => {
                    return (
                      <Tag key={categ} color={'white'} bg={'blue.400'} mr={2} mb={1}>
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

export default BasicDetailsReviewbasicDetailTable;
