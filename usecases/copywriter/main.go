package copywriter

import (
	"context"
	"fmt"
	"time"

	"github.com/nullswan/nomi/internal/chat"
	"github.com/nullswan/nomi/internal/tools"
)

// TODO(nullswan): Pull preferences / Profiles / Previous projects
// TODO(nullswan): Add Memory, Storage, Tools
// TODO(nullswan): Generate multiple content examples
// TODO(nullswan): Add ref check agent

func OnStart(
	ctx context.Context,
	selector tools.Selector,
	logger tools.Logger,
	inputHandler tools.InputHandler,
	textToJSONBackend tools.TextToJSONBackend,
	conversation chat.Conversation, // TOOD(nullswan): Should be like a project
) error {
	project := fmt.Sprintf(
		"project-copywriting-%s",
		time.Now().Format("2006-01-02"),
	)

	goalsAgent := NewGoalsAgent(
		textToJSONBackend,
		inputHandler,
		logger,
		selector,
	)
	ideasAgent := NewIdeasAgent(
		logger,
		textToJSONBackend,
		goalsAgent,
		inputHandler,
	)
	headlineAgent := NewHeadlineAgent(
		logger,
		textToJSONBackend,
		selector,
		inputHandler,
		goalsAgent,
		ideasAgent,
	)
	contentPlanAgent := NewOutlineAgent(
		logger,
		textToJSONBackend,
		selector,
		inputHandler,
		goalsAgent,
		ideasAgent,
		headlineAgent,
	)
	exportAgent := NewExportAgent(
		logger,
		project,
	)
	redactAgent := NewRedactAgent(
		goalsAgent,
		ideasAgent,
		headlineAgent,
		contentPlanAgent,
		exportAgent,
		logger,
		inputHandler,
		textToJSONBackend,
		selector,
	)

	err := goalsAgent.OnStart(ctx, conversation)
	if err != nil {
		return fmt.Errorf("error starting goals agent: %w", err)
	}
	conversation = conversation.Reset()

	err = ideasAgent.OnStart(ctx, conversation)
	if err != nil {
		return fmt.Errorf("error starting ideas agent: %w", err)
	}
	conversation = conversation.Reset()

	err = headlineAgent.OnStart(ctx, conversation)
	if err != nil {
		return fmt.Errorf("error starting headline agent: %w", err)
	}

	conversation = conversation.Reset()
	err = contentPlanAgent.OnStart(ctx, conversation)
	if err != nil {
		return fmt.Errorf("error starting content plan agent: %w", err)
	}

	conversation = conversation.Reset()
	err = redactAgent.OnStart(
		ctx,
		conversation,
	)
	if err != nil {
		return fmt.Errorf("error starting redact agent: %w", err)
	}

	return nil
}
