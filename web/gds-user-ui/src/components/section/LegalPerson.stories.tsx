import { Meta, Story } from "@storybook/react";
import LegalPerson from "./LegalPerson";

type LegalPersonProps = {};

export default {
  title: "components/LegalPerson",
  component: LegalPerson,
} as Meta<LegalPersonProps>;

const Template: Story<LegalPersonProps> = (args) => <LegalPerson {...args} />;

export const Default = Template.bind({});
