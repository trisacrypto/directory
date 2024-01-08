import { Meta, Story } from "@storybook/react";
import NationalIdentification from ".";

type NationalIdentificationProps = {};

export default {
  title: "components/NationalIdentification",
  component: NationalIdentification,
} as Meta<NationalIdentificationProps>;

const Template: Story<NationalIdentificationProps> = (args) => (
  <NationalIdentification {...args} />
);

export const Default = Template.bind({});
Default.args = {};
