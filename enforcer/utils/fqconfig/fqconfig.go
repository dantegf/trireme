package fqconfig

import "strconv"

// FilterQueue captures all the configuration parameters of the NFQUEUEs
type FilterQueue struct {
	// queueSeparation specifies if we should use separate queues per packet type
	queueSeparation bool
	// MarkValue is the default mark to set in packets in the RAW chain
	markValue int
	// Network Queue is the queue number of the base queue for network packets
	networkQueue uint16
	// NumberOfApplicationQueues is the number of queues that must be allocated
	numberOfApplicationQueues uint16
	// NumberOfNetworkQueues is the number of network queues allocated
	numberOfNetworkQueues uint16
	// ApplicationQueue is the queue number of the first application queue
	applicationQueue uint16
	// ApplicationQueueSize is the size of the application queue
	applicationQueueSize uint32
	// NetworkQueueSize is the size of the network queue
	networkQueueSize uint32

	// Strings for programming NFQ
	networkQueuesSynStr     string
	networkQueuesAckStr     string
	networkQueuesSvcStr     string
	applicationQueuesSynStr string
	applicationQueuesAckStr string
	applicationQueuesSvcStr string
}

// FilterQueueImpl is the interface for filter queue configs
type FilterQueueImpl interface {

	// GetMarkValue returns a mark value to be used by iptables action
	GetMarkValue() int

	// GetNetworkQueueStart returns start of network queues to be used by iptables action
	GetNetworkQueueStart() uint16

	// GetNumNetworkQueues returns number of network queues to be used by iptables action
	GetNumNetworkQueues() uint16

	// GetNetworkQueueSize returns size of network queues to be used by iptables action
	GetNetworkQueueSize() uint32

	// GetApplicationQueueStart returns start of application queues to be used by iptables action
	GetApplicationQueueStart() uint16

	// GetNumApplicationQueues returns number of application queues to be used by iptables action
	GetNumApplicationQueues() uint16

	// GetApplicationQueueSize returns size of application queues to be used by iptables action
	GetApplicationQueueSize() uint32

	// GetNetworkQueueSynStr returns a queue id string to be used by iptables action
	GetNetworkQueueSynStr() string

	// GetNetworkQueueAckStr returns a queue id string to be used by iptables action
	GetNetworkQueueAckStr() string

	// GetNetworkQueueSvcStr returns a queue id string to be used by iptables action
	GetNetworkQueueSvcStr() string

	// GetApplicationQueueSynStr returns a queue id string to be used by iptables action
	GetApplicationQueueSynStr() string

	// GetApplicationQueueAckStr returns a queue id string to be used by iptables action
	GetApplicationQueueAckStr() string

	// GetApplicationQueueSvcStr returns a queue id string to be used by iptables action
	GetApplicationQueueSvcStr() string
}

// NewFilterQueueWithDefaults return a default filter queue config
func NewFilterQueueWithDefaults() FilterQueueImpl {
	return NewFilterQueue(
		DefaultQueueSeperation,
		DefaultMarkValue,
		DefaultQueueStart,
		DefaultNumberOfQueues,
		DefaultNumberOfQueues,
		DefaultQueueSize,
		DefaultQueueSize,
	)
}

// NewFilterQueue returns an instance of FilterQueue
func NewFilterQueue(queueSeparation bool, MarkValue int, QueueStart, NumberOfNetworkQueues, NumberOfApplicationQueues uint16, NetworkQueueSize, ApplicationQueueSize uint32) FilterQueueImpl {

	fq := &FilterQueue{
		queueSeparation:      queueSeparation,
		markValue:            MarkValue,
		networkQueueSize:     NetworkQueueSize,
		applicationQueueSize: ApplicationQueueSize,
	}

	if queueSeparation {

		fq.applicationQueue = QueueStart
		fq.applicationQueuesSynStr = strconv.Itoa(int(fq.applicationQueue)) + ":" + strconv.Itoa(int(fq.applicationQueue+NumberOfApplicationQueues-1))
		fq.applicationQueuesAckStr = strconv.Itoa(int(fq.applicationQueue+1*NumberOfApplicationQueues)) + ":" + strconv.Itoa(int(fq.applicationQueue+2*NumberOfApplicationQueues-1))
		fq.applicationQueuesSvcStr = strconv.Itoa(int(fq.applicationQueue+2*NumberOfApplicationQueues)) + ":" + strconv.Itoa(int(fq.applicationQueue+3*NumberOfApplicationQueues-1))
		fq.numberOfApplicationQueues = NumberOfApplicationQueues * 3

		fq.networkQueue = QueueStart + fq.numberOfApplicationQueues
		fq.networkQueuesSynStr = strconv.Itoa(int(fq.networkQueue)) + ":" + strconv.Itoa(int(fq.networkQueue+NumberOfNetworkQueues-1))
		fq.networkQueuesAckStr = strconv.Itoa(int(fq.networkQueue+1*NumberOfNetworkQueues)) + ":" + strconv.Itoa(int(fq.networkQueue+2*NumberOfNetworkQueues-1))
		fq.networkQueuesSvcStr = strconv.Itoa(int(fq.networkQueue+2*NumberOfNetworkQueues)) + ":" + strconv.Itoa(int(fq.networkQueue+3*NumberOfNetworkQueues-1))
		fq.numberOfNetworkQueues = NumberOfNetworkQueues * 3
	} else {

		fq.applicationQueue = QueueStart
		fq.applicationQueuesSynStr = strconv.Itoa(int(fq.applicationQueue)) + ":" + strconv.Itoa(int(fq.applicationQueue+NumberOfApplicationQueues-1))
		fq.applicationQueuesAckStr = fq.applicationQueuesSynStr
		fq.applicationQueuesSvcStr = fq.applicationQueuesSynStr
		fq.numberOfApplicationQueues = NumberOfApplicationQueues

		fq.networkQueue = QueueStart + fq.numberOfApplicationQueues
		fq.networkQueuesSynStr = strconv.Itoa(int(fq.networkQueue)) + ":" + strconv.Itoa(int(fq.networkQueue+NumberOfNetworkQueues-1))
		fq.networkQueuesAckStr = fq.networkQueuesSynStr
		fq.networkQueuesSvcStr = fq.networkQueuesSynStr
		fq.numberOfNetworkQueues = NumberOfNetworkQueues
	}
	return fq
}

// GetMarkValue returns a mark value to be used by iptables action
func (f *FilterQueue) GetMarkValue() int {
	return f.markValue
}

// GetNetworkQueueStart returns start of network queues to be used by iptables action
func (f *FilterQueue) GetNetworkQueueStart() uint16 {
	return f.networkQueue
}

// GetNumNetworkQueues returns number of network queues to be used by iptables action
func (f *FilterQueue) GetNumNetworkQueues() uint16 {
	return f.numberOfNetworkQueues
}

// GetNetworkQueueSize returns size of network queues to be used by iptables action
func (f *FilterQueue) GetNetworkQueueSize() uint32 {
	return f.networkQueueSize
}

// GetApplicationQueueStart returns start of application queues to be used by iptables action
func (f *FilterQueue) GetApplicationQueueStart() uint16 {
	return f.applicationQueue
}

// GetNumApplicationQueues returns number of application queues to be used by iptables action
func (f *FilterQueue) GetNumApplicationQueues() uint16 {
	return f.numberOfApplicationQueues
}

// GetApplicationQueueSize returns size of application queues to be used by iptables action
func (f *FilterQueue) GetApplicationQueueSize() uint32 {
	return f.applicationQueueSize
}

// GetNetworkQueueSynStr returns a queue id string to be used by iptables action
func (f *FilterQueue) GetNetworkQueueSynStr() string {
	return f.networkQueuesSynStr
}

// GetNetworkQueueAckStr returns a queue id string to be used by iptables action
func (f *FilterQueue) GetNetworkQueueAckStr() string {
	if f.queueSeparation {
		return f.networkQueuesAckStr
	}
	return f.networkQueuesSynStr
}

// GetNetworkQueueSvcStr returns a queue id string to be used by iptables action
func (f *FilterQueue) GetNetworkQueueSvcStr() string {
	if f.queueSeparation {
		return f.networkQueuesSvcStr
	}
	return f.networkQueuesSynStr
}

// GetApplicationQueueSynStr returns a queue id string to be used by iptables action
func (f *FilterQueue) GetApplicationQueueSynStr() string {
	return f.applicationQueuesSynStr
}

// GetApplicationQueueAckStr returns a queue id string to be used by iptables action
func (f *FilterQueue) GetApplicationQueueAckStr() string {
	if f.queueSeparation {
		return f.applicationQueuesAckStr
	}
	return f.applicationQueuesSynStr
}

// GetApplicationQueueSvcStr returns a queue id string to be used by iptables action
func (f *FilterQueue) GetApplicationQueueSvcStr() string {
	if f.queueSeparation {
		return f.applicationQueuesSvcStr
	}
	return f.applicationQueuesSynStr
}

// Default parameters for the NFQUEUE configuration. Parameters can be
// changed after an isolator has been created and before its started.
// Change in parameters after the isolator is started has no effect
const (
	// DefaultQueueSeperation specifies if we should use separate queues for packet types
	DefaultQueueSeperation = true
	// DefaultNumberOfQueues  is the default number of queues used in NFQUEUE
	DefaultNumberOfQueues = 4
	// DefaultQueueStart represents the queue number to start
	DefaultQueueStart = 0
	// DefaultQueueSize is the size of the queues
	DefaultQueueSize = 500
	// DefaultMarkValue is the default Mark for packets in the raw chain
	DefaultMarkValue = 0x1111
)