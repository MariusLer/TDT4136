# multiAgents.py
# --------------
# Licensing Information:  You are free to use or extend these projects for
# educational purposes provided that (1) you do not distribute or publish
# solutions, (2) you retain this notice, and (3) you provide clear
# attribution to UC Berkeley, including a link to http://ai.berkeley.edu.
#
# Attribution Information: The Pacman AI projects were developed at UC Berkeley.
# The core projects and autograders were primarily created by John DeNero
# (denero@cs.berkeley.edu) and Dan Klein (klein@cs.berkeley.edu).
# Student side autograding was added by Brad Miller, Nick Hay, and
# Pieter Abbeel (pabbeel@cs.berkeley.edu).


from util import manhattanDistance
from game import Directions
import random, util

from game import Agent

class ReflexAgent(Agent):
    """
      A reflex agent chooses an action at each choice point by examining
      its alternatives via a state evaluation function.

      The code below is provided as a guide.  You are welcome to change
      it in any way you see fit, so long as you don't touch our method
      headers.
    """


    def getAction(self, gameState):
        """
        You do not need to change this method, but you're welcome to.

        getAction chooses among the best options according to the evaluation function.

        Just like in the previous project, getAction takes a gameState and returns
        some Directions.X for some X in the set {North, South, West, East, Stop}
        """
        # Collect legal moves and successor states
        legalMoves = gameState.getLegalActions()

        # Choose one of the best actions
        scores = [self.evaluationFunction(gameState, action) for action in legalMoves]
        bestScore = max(scores)
        bestIndices = [index for index in range(len(scores)) if scores[index] == bestScore]
        chosenIndex = random.choice(bestIndices) # Pick randomly among the best

        "Add more of your code here if you want to"

        return legalMoves[chosenIndex]

    def evaluationFunction(self, currentgameState, action):
        """
        Design a better evaluation function here.

        The evaluation function takes in the current and proposed successor
        gameStates (pacman.py) and returns a number, where higher numbers are better.

        The code below extracts some useful information from the state, like the
        remaining food (newFood) and Pacman position after moving (newPos).
        newScaredTimes holds the number of moves that each ghost will remain
        scared because of Pacman having eaten a power pellet.

        Print out these variables to see what you're getting, then combine them
        to create a masterful evaluation function.
        """
        # Useful information you can extract from a gameState (pacman.py)
        successorgameState = currentgameState.generatePacmanSuccessor(action)
        newPos = successorgameState.getPacmanPosition()
        newFood = successorgameState.getFood()
        newGhostStates = successorgameState.getGhostStates()
        newScaredTimes = [ghostState.scaredTimer for ghostState in newGhostStates]

        "*** YOUR CODE HERE ***"
        return successorgameState.getScore()

def scoreEvaluationFunction(currentgameState):
    """
      This default evaluation function just returns the score of the state.
      The score is the same one displayed in the Pacman GUI.

      This evaluation function is meant for use with adversarial search agents
      (not reflex agents).
    """
    return currentgameState.getScore()

class MultiAgentSearchAgent(Agent):
    """
      This class provides some common elements to all of your
      multi-agent searchers.  Any methods defined here will be available
      to the MinimaxPacmanAgent, AlphaBetaPacmanAgent & ExpectimaxPacmanAgent.

      You *do not* need to make any changes here, but you can if you want to
      add functionality to all your adversarial search agents.  Please do not
      remove anything, however.

      Note: this is an abstract class: one that should not be instantiated.  It's
      only partially specified, and designed to be extended.  Agent (game.py)
      is another abstract class.
    """

    def __init__(self, evalFn = 'scoreEvaluationFunction', depth = '2'):
        self.index = 0 # Pacman is always agent index 0
        self.evaluationFunction = util.lookup(evalFn, globals())
        self.depth = int(depth)

class MinimaxAgent(MultiAgentSearchAgent):
    """
      Your minimax agent (question 2)
    """

    def getAction(self, gameState):
        """
          Returns the minimax action from the current gameState using self.depth
          and self.evaluationFunction.

          Here are some method calls that might be useful when implementing minimax.

          gameState.getLegalActions(agentIndex):
            Returns a list of legal actions for an agent
            agentIndex=0 means Pacman, ghosts are >= 1

          gameState.generateSuccessor(agentIndex, action):
            Returns the successor game state after an agent takes an action

          gameState.getNumAgents():
            Returns the total number of agents in the game
        """
        "*** YOUR CODE HERE ***"
        bestValueAction = self.maximize(gameState, 0)
        return bestValueAction[1]

    # I used a touple to hold the bestvalue and bestaction found in every step. Althoug the best action is actually only needed for the last step when we return action

    def maximize(self, gameState, currentDepth):   #// want to maximize own utility
            if currentDepth == self.depth or gameState.isWin() or gameState.isLose(): # Reached terminal nodes
                return (self.evaluationFunction(gameState),None)
            bestValueAction = (float("-inf"),"")
            actions = gameState.getLegalActions(0)
            for action in actions:
                successor = gameState.generateSuccessor(0,action)
                successorValue = self.minimize(successor,currentDepth,1) # Ghosts mimimize utility
                if successorValue[0] > bestValueAction[0]:  # Found a better path
                    bestValueAction = (successorValue[0],action)
            return bestValueAction # return the best possible value and action found

    def minimize(self,gameState, currentDepth, agentIndex): # ghosts time to minimize utility
            if gameState.isWin() or gameState.isLose(): # Reached terminal nodes
                return (self.evaluationFunction(gameState),None)
            bestValueAction = (float("inf"),"")
            numbAgents = gameState.getNumAgents()
            actions = gameState.getLegalActions(agentIndex)
            for action in actions:
                successor = gameState.generateSuccessor(agentIndex,action)
                if agentIndex == numbAgents-1:  # Check if we are the last ghost to run minimize, if so run maximize on pacman again
                    successorValue = self.maximize(successor, currentDepth+1)
                else:
                    successorValue = self.minimize(successor, currentDepth, agentIndex+1) # There are more ghosts, so we have to run minimize again
                if successorValue[0] < bestValueAction[0]: # Value that minimizes found
                    bestValueAction = (successorValue[0],action)
            return bestValueAction # Return lowest value/action found

class AlphaBetaAgent(MultiAgentSearchAgent):
    """
      Your minimax agent with alpha-beta pruning (question 3)
    """

    def getAction(self, gameState):
        """
          Returns the minimax action using self.depth and self.evaluationFunction
        """
        "*** YOUR CODE HERE ***"
        bestValueAction = self.maximize(gameState, 0, float("-inf"), float("inf"))
        return bestValueAction[1]

    def maximize(self, gameState, currentDepth, alpha, beta):   #// want to maximize own utility
            if currentDepth == self.depth or gameState.isWin() or gameState.isLose(): # Reached terminal nodes or lost
                return (self.evaluationFunction(gameState),None)  # Return none for no action here
            bestValueAction = (float("-inf"),"")   # initialize best value and action
            actions = gameState.getLegalActions(0)
            for action in actions:
                successor = gameState.generateSuccessor(0,action)
                successorValue = self.minimize(successor,currentDepth,1,alpha,beta) # Ghosts mimimize utility of pacman
                if successorValue[0] > bestValueAction[0]:
                    bestValueAction = (successorValue[0],action)
                if bestValueAction[0] > beta:
                    return bestValueAction
                alpha = max(alpha, bestValueAction[0])
            return bestValueAction

    def minimize(self, gameState, currentDepth, agentIndex, alpha, beta): # ghosts time to minimize utility
            if gameState.isWin() or gameState.isLose(): # Reached terminal nodes
                return (self.evaluationFunction(gameState),None)
            bestValueAction = (float("inf"),"")
            numbAgents = gameState.getNumAgents()
            actions = gameState.getLegalActions(agentIndex)
            for action in actions:
                successor = gameState.generateSuccessor(agentIndex,action)
                if agentIndex == numbAgents-1:
                    successorValue = self.maximize(successor, currentDepth+1,alpha,beta)
                else:
                    successorValue = self.minimize(successor, currentDepth, agentIndex+1,alpha,beta)
                if successorValue[0] < bestValueAction[0]:
                    bestValueAction = (successorValue[0],action)
                if bestValueAction[0] < alpha:
                    return bestValueAction
                beta = min(beta, bestValueAction[0])
            return bestValueAction


class ExpectimaxAgent(MultiAgentSearchAgent):
    """
      Your expectimax agent (question 4)
    """

    def getAction(self, gameState):
        """
          Returns the expectimax action using self.depth and self.evaluationFunction

          All ghosts should be modeled as choosing uniformly at random from their
          legal moves.
        """
        "*** YOUR CODE HERE ***"
        util.raiseNotDefined()

def betterEvaluationFunction(currentgameState):
    """
      Your extreme ghost-hunting, pellet-nabbing, food-gobbling, unstoppable
      evaluation function (question 5).

      DESCRIPTION: <write something here so we know what you did>
    """
    "*** YOUR CODE HERE ***"
    util.raiseNotDefined()

# Abbreviation
better = betterEvaluationFunction
