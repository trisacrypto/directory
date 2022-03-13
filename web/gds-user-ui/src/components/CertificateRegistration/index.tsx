import React from "react";
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
} from "@chakra-ui/react";
import CertificateRegistrationRow from "components/Tables/CertificateRegistrationRow";
import Card from "components/ui/Card";
const defaultRowData = [
  {
    section: "1",
    name: "Business Details",
    description: "Website, incorporation Date, VASP Category",
    status: null,
  },
  {
    section: "2",
    name: "Legal Person",
    description: "Name, Addresss, Country, National Identifier",
    status: null,
  },
  {
    section: "3",
    name: "Contacts",
    description: "Compliance, Technical, Admininstrative, Billing",
    status: null,
  },
  {
    section: "4",
    name: "TRISA Implementation",
    description: "TRISA endpoint for communication",
    status: null,
  },
  {
    section: "5",
    name: "TRIXO Questionnaire",
    description: "CDD and data protection policies",
    status: null,
  },
  {
    section: "6",
    name: "Review & Submit",
    description: "Final review and form submission",
    status: null,
  },
];
const CertificateRegistration = ({ title, data }: any) => {
  const textColor = useColorModeValue("#858585", "white");
  return (
    <Card>
      <Stack p={4} mb={5}>
        <Heading fontSize="20px" fontWeight="bold" pb=".5rem">
          Certificate Registration Process
        </Heading>
      </Stack>
      <Stack px={"60px"}>
        <Table
          color={textColor}
          width={"100%"}
          sx={{
            borderCollapse: "separate",
            borderSpacing: "0 10px",
            Th: {
              textTransform: "capitalize",
              color: "#858585",
              fontWeight: "bold",
              borderBottom: "none",
              fontSize: "0.9rem",
              textAlign: "center",
            },
          }}
        >
          <Thead>
            <Tr>
              <Th>Section</Th>
              <Th>Name</Th>
              <Th>Description</Th>
              <Th>Status</Th>
              <Th>Action</Th>
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
