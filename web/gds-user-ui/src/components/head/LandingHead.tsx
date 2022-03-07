import React from "react";
import { Flex, Heading, Stack, Text } from "@chakra-ui/react";

interface LandingHeaderProps {
  title: string;
  description?: string;
}
// we should add props to the LandingHead component to allow it to update content dynamically
const LandingHead: React.FC<any> = ({
  title,
  description,
}: LandingHeaderProps): any => {
  return (
    <Flex
      bgGradient="linear(270deg,#24a9df,#1aebb4)"
      color="white"
      width="100%"
      justifyContent="center"
      alignItems="center"
      direction="column"
      padding={4}
    >
      <Stack
        textAlign={"center"}
        color="white"
        align={"center"}
        spacing={{ base: 3 }}
        py={{ base: 5 }}
      >
        <Heading
          fontWeight={600}
          fontSize={{ base: "2xl", sm: "2xl", md: "4xl" }}
          lineHeight={"80%"}
        >
          {title || "TRISA Global Directory Service"}
        </Heading>
        {description ? (
          <Text maxW={"2xl"}>{description}</Text>
        ) : (
          <Text fontSize={"xl"}>
            Become Travel Rule compliant. <br />
            Apply to Become a TRISA certified Virtual Asset Service Provider.'
          </Text>
        )}
      </Stack>
    </Flex>
  );
};

export default LandingHead;
