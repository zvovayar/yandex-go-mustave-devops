package grpcs

import (
	"context"
	"fmt"

	"github.com/zvovayar/yandex-go-mustave-devops/internal/proto"
	inst "github.com/zvovayar/yandex-go-mustave-devops/internal/storage"
)

type UsersServer struct {
	// нужно встраивать тип pb.Unimplemented<TypeName>
	// для совместимости с будущими версиями
	proto.UnimplementedUsersServer
}

// GetMetric return Metric value
func (us *UsersServer) GetMetric(ctx context.Context, m *proto.Metric) (*proto.Metric, error) {

	if m.MType == "gauge" {
		m.Value = float64(inst.StoreMonitor.GetGMvalue(m.ID))
	}
	if m.MType == "counter" {
		m.Delta = int64(inst.StoreMonitor.GetCMvalue(m.ID))
	}
	return nil, nil
}

// UpdateMetric save Metric value
func (us *UsersServer) UpdateMetric(ctx context.Context, m *proto.Metric) (*proto.Result, error) {

	//
	// сохранять значение метрики
	//
	if m.MType == "gauge" {

		inst.Sugar.Infof("m.Value=%f", m.Value)
		gmname := m.ID
		gm := inst.Gauge(m.Value)

		inst.StoreMonitor.SetGMvalue(gmname, inst.Gauge(gm))

		inst.Sugar.Infof("Store %v = %f", gmname, gm)

		return &proto.Result{
			Code: 0,
			Text: fmt.Sprintf("Store type=%s %v = %f", m.MType, m.ID, m.Value),
		}, nil

	} else if m.MType == "counter" {

		inst.Sugar.Infof("m.Delta=%d", m.Delta)
		cmname := m.ID
		cm := inst.Counter(m.Delta)

		inst.StoreMonitor.SetCMvalue(cmname, inst.Counter(cm))

		inst.Sugar.Infof("Store %v = %d", cmname, cm)
	}
	return &proto.Result{
		Code: 0,
		Text: fmt.Sprintf("Store type=%s %v = %d", m.MType, m.ID, m.Delta),
	}, nil

}

// UpdateBatchMetrics save array of Metrics
func (us *UsersServer) UpdateBatchMetrics(ctx context.Context, m *proto.BatchMetrics) (*proto.Result, error) {

	var mbatch []inst.Metrics

	// convert from protobuf to []inst.Metrics
	for i := 0; i < int(m.Count); i++ {

		var mbuf = inst.Metrics{
			ID:    m.Metrics[i].ID,
			MType: m.Metrics[i].MType,
			Delta: &m.Metrics[i].Delta,
			Value: &m.Metrics[i].Value,
			Hash:  m.Metrics[i].Hash,
		}

		mbatch = append(mbatch, mbuf)
	}

	// save batch
	if err := inst.StoreMonitor.SaveBatch(context.Background(), mbatch); err != nil {
		inst.Sugar.Infow(err.Error())
		return &proto.Result{
			Code: 1,
			Text: "can not save batch " + err.Error(),
		}, err
	}

	return &proto.Result{
		Code: 0,
		Text: fmt.Sprintf("save batch with %v values", m.Count),
	}, nil
}
