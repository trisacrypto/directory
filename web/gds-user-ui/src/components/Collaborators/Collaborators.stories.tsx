import { Meta, Story } from "@storybook/react";
import CollaboratorsSection from ".";

type CollaboratorsSectionProps = {};

export default {
  title: "components/CollaboratorsSection",
  component: CollaboratorsSection,
} as Meta<CollaboratorsSectionProps>;

const Template: Story<CollaboratorsSectionProps> = (args) => (
  <CollaboratorsSection {...args} />
);

export const Standard = Template.bind({});
Standard.args = {};
