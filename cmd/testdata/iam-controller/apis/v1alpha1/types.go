package v1alpha1

// RoleSpec defines the desired state of Role.
type RoleSpec struct {
	// AssumeRolePolicyDocument is the trust policy for the IAM role.
	// +kubebuilder:validation:Required
	AssumeRolePolicyDocument *string `json:"assumeRolePolicyDocument,omitempty"`

	// Description is a description of the role.
	Description *string `json:"description,omitempty"`

	// MaxSessionDuration is the max session duration in seconds.
	MaxSessionDuration *int64 `json:"maxSessionDuration,omitempty"`

	// PermissionsBoundaryPolicyARN is the ARN of the policy for permissions boundary.
	PermissionsBoundaryPolicyARN *string `json:"permissionsBoundaryPolicyARN,omitempty"`

	// RoleName is the name of the IAM role.
	RoleName *string `json:"roleName,omitempty"`

	// InlinePolicyDocument contains the inline policy.
	InlinePolicyDocument *string `json:"inlinePolicyDocument,omitempty"`

	// Tags contains the tags for the role.
	Tags []*Tag `json:"tags,omitempty"`
}

// RoleStatus defines the observed state of Role.
type RoleStatus struct {
	// ARN is the Amazon Resource Name for the role.
	ARN *string `json:"arn,omitempty"`

	// RoleID is the stable unique identifier for the role.
	RoleID *string `json:"roleID,omitempty"`
}

// PolicySpec defines the desired state of Policy.
type PolicySpec struct {
	// PolicyDocument is the JSON policy document.
	// +kubebuilder:validation:Required
	PolicyDocument *string `json:"policyDocument,omitempty"`

	// PolicyName is the friendly name for the policy.
	PolicyName *string `json:"policyName,omitempty"`

	// Description is a description for the policy.
	Description *string `json:"description,omitempty"`
}

// Tag represents an IAM tag.
type Tag struct {
	Key   *string `json:"key,omitempty"`
	Value *string `json:"value,omitempty"`
}
